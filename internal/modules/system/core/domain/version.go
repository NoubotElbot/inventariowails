package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a semantic version reduced to its numeric core (X.Y.Z).
type Version struct {
	Major, Minor, Patch int
}

// ParseVersion parses "v1.2.3" or "1.2.3". Pre-release and build suffixes
// ("-rc.1", "+abc") are ignored. Non-release builds such as "dev" fail.
func ParseVersion(s string) (Version, error) {
	raw := strings.TrimPrefix(strings.TrimSpace(s), "v")
	if i := strings.IndexAny(raw, "-+"); i >= 0 {
		raw = raw[:i]
	}

	parts := strings.Split(raw, ".")
	if len(parts) == 0 || len(parts) > 3 {
		return Version{}, fmt.Errorf("versión inválida: %q", s)
	}

	nums := [3]int{}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return Version{}, fmt.Errorf("versión inválida: %q", s)
		}
		nums[i] = n
	}

	return Version{Major: nums[0], Minor: nums[1], Patch: nums[2]}, nil
}

// Compare returns -1, 0 or 1 if v is lower, equal or greater than o.
func (v Version) Compare(o Version) int {
	for _, d := range [...]int{v.Major - o.Major, v.Minor - o.Minor, v.Patch - o.Patch} {
		if d < 0 {
			return -1
		}
		if d > 0 {
			return 1
		}
	}
	return 0
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
