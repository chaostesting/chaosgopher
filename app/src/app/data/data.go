package data

import "errors"
import "fmt"
import "math/rand"
import "time"

var (
	ErrNoSources    = errors.New("there is no data source configured.")
	ErrNoGoodSource = errors.New("unable to establish a good connection to any data source, check your servers.")
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Interface interface {
	Get(result interface{}, query string, args ...interface{}) error
	Select(result interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) error
	Close() error
}

type Source struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
	Weight int    `json:"weight"`
}

// SelectSource() selects a random Source from the given array while taking the weight of
// each source into account.
func SelectSource(sources []*Source) *Source {
	if len(sources) <= 0 {
		return nil
	}

	sums := make([]int, len(sources))
	for i, source := range sources {
		if i == 0 {
			sums[i] = source.Weight
			continue
		}

		sums[i] = sums[i-1] + source.Weight
	}

	weight := rand.Intn(sums[len(sums)-1])
	for i := len(sums) - 1; i >= 0; i-- {
		if sums[i] < weight {
			return sources[i]
		}
	}

	return sources[0]
}

func (s *Source) String() string {
	return fmt.Sprintf("%s (%d)", s.Name, s.Weight)
}
