package routes

import (
	"cmp"
	"fmt"
	"net"
	"net/netip"
	"slices"
	"strings"
)

type Route struct {
	Prefix  netip.Prefix
	Comment string
}

type RouteList struct {
	Routes   []Route
	Count    int
	Deleted  int
	IsSorted bool
}

// NewRouteList Создает новый список маршрутов из переданного среза маршрутов. Если срез пуст, возвращает nil.
func NewRouteList(routes []Route) *RouteList {
	if len(routes) == 0 {
		return nil
	}
	return &RouteList{
		Routes:   routes,
		Count:    len(routes),
		Deleted:  0,
		IsSorted: false,
	}
}

func ParseRouteFromLine(line, source string) (Route, error) {
	comment := source
	line = strings.ToLower(strings.TrimSpace(line))
	if !strings.HasPrefix(line, "route add") {
		return Route{}, fmt.Errorf("Строка не начинается с route add")
	}
	partsByRem := strings.Split(line, "& rem")
	if len(partsByRem) == 2 {
		comment = strings.TrimSpace(partsByRem[1])
	}
	parts := strings.Fields(partsByRem[0])
	if len(parts) < 5 {
		return Route{}, fmt.Errorf("Недостаточно информации в строке")
	}
	addr, err := netip.ParseAddr(parts[2])
	if err != nil {
		return Route{}, fmt.Errorf("Неверный IP-адрес: %s", parts[2])
	}
	maskIP := net.ParseIP(parts[4]).To4()
	if maskIP == nil {
		return Route{}, fmt.Errorf("Неверный ip адрес маски: %s", parts[4])
	}
	mask, _ := net.IPMask(maskIP).Size()
	return Route{
		Prefix:  netip.PrefixFrom(addr, mask),
		Comment: comment,
	}, nil
}

func (rt *RouteList) Sort() {
	slices.SortFunc(rt.Routes, func(a, b Route) int {
		if a.Prefix == b.Prefix {
			return cmp.Compare(a.Comment, b.Comment)
		}
		aNet := a.Prefix.Masked()
		bNet := b.Prefix.Masked()
		if c := aNet.Addr().Compare(bNet.Addr()); c != 0 {
			return c
		}
		return cmp.Compare(a.Prefix.Bits(), b.Prefix.Bits())
	})
	rt.IsSorted = true
}

func (rt *RouteList) UniqueRoutes() (RouteList, error) {
	if !rt.IsSorted {
		return RouteList{}, fmt.Errorf("Маршруты не отсортированы")
	}
	var unique []Route
	var deleted int
	count := len(rt.Routes)
	unique = append(unique, rt.Routes[0])
	for i := 1; i < count; i++ {
		parent := unique[len(unique)-1].Prefix.Masked()
		sub := rt.Routes[i].Prefix.Masked()
		if parent.Bits() <= sub.Bits() && parent.Contains(sub.Addr()) {
			deleted++
			continue
		} else {
			unique = append(unique, rt.Routes[i])
		}
	}

	return RouteList{
		Routes:   unique,
		Count:    len(unique),
		Deleted:  deleted,
		IsSorted: true,
	}, nil
}

func LineToBat(route Route) string {
	prefix := route.Prefix.Masked()
	ip := prefix.Addr()
	mask := net.CIDRMask(prefix.Bits(), 32)
	maskIP := net.IP(mask)
	return fmt.Sprintf("route ADD %s MASK %s 0.0.0.0 & rem %s", ip.String(), maskIP.String(), route.Comment)
}
