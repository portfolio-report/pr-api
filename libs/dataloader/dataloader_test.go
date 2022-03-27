package dataloader_test

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/portfolio-report/pr-api/libs/dataloader"
	"github.com/stretchr/testify/require"
)

type testValue struct {
	ID   string
	Name string
}

// GetMockedLoader returns a dataloader and pointer to the list of
// all calls with their keys to the mocked fetch function
func GetMockedLoader(max int) (*dataloader.Dataloader[string, *testValue], *[][]string) {
	var mu sync.Mutex
	var fetchCalls [][]string

	loader := dataloader.New(dataloader.Config[string, *testValue]{
		Fetch: func(keys []string) ([]*testValue, []error) {
			mu.Lock()
			fetchCalls = append(fetchCalls, keys)
			mu.Unlock()

			values := make([]*testValue, len(keys))
			errors := make([]error, len(keys))

			for i, key := range keys {
				if strings.HasPrefix(key, "E") {
					errors[i] = fmt.Errorf("not found " + key)
				} else {
					values[i] = &testValue{ID: key, Name: "value " + key}
				}
			}
			return values, errors
		},
		MaxSize: max,
	},
	)
	return loader, &fetchCalls
}

func TestDataloader(t *testing.T) {
	t.Run("successful Load", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		value, err := dl.Load("1")
		require.Nil(t, err)
		require.Equal(t, "1", value.ID)
	})

	t.Run("fetch function called only once when cached", func(t *testing.T) {
		t.Parallel()
		dl, fetchCallsPtr := GetMockedLoader(0)
		for i := 0; i < 2; i++ {
			_, err := dl.Load("E1")
			require.Error(t, err)
			fetchCalls := *fetchCallsPtr
			require.Len(t, fetchCalls, 1)
			require.Len(t, fetchCalls[0], 1)
		}
		for i := 0; i < 2; i++ {
			val, err := dl.Load("U1")
			require.Equal(t, "U1", val.ID)
			require.NoError(t, err)
			fetchCalls := *fetchCallsPtr
			require.Len(t, fetchCalls, 2)
			require.Len(t, fetchCalls[1], 1)
		}
	})

	t.Run("failed Load", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		value, err := dl.Load("E1")
		require.Error(t, err)
		require.Nil(t, value)
	})

	t.Run("LoadMany", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		u, err := dl.LoadMany([]string{"U1", "U2", "E3", "E4", "U5"})
		require.Equal(t, u[0].ID, "U1")
		require.Equal(t, u[1].ID, "U2")
		require.Error(t, err[2])
		require.Error(t, err[3])
		require.Equal(t, u[4].ID, "U5")
	})

	t.Run("LoadThunk does not contain race conditions", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		future := dl.LoadThunk("1")
		go future()
		go future()
	})

	t.Run("LoadThunk", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		thunk1 := dl.LoadThunk("U1")
		thunk2 := dl.LoadThunk("E1")

		u1, err1 := thunk1()
		require.NoError(t, err1)
		require.Equal(t, "value U1", u1.Name)

		u2, err2 := thunk2()
		require.Error(t, err2)
		require.Nil(t, u2)
	})

	t.Run("LoadMany returns errors", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		_, err := dl.LoadMany([]string{"E1", "E2", "E3"})
		require.Len(t, err, 3)
	})

	t.Run("LoadMany returns len(errors) == len(keys)", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(3)
		_, errs := dl.LoadMany([]string{"E1", "U2", "U3"})
		require.Len(t, errs, 3)

		require.Error(t, errs[0])
		require.Nil(t, errs[1])
		require.Nil(t, errs[2])
	})

	t.Run("LoadMany returns nil []error when no errors occurred", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		_, errs := dl.LoadMany([]string{"1", "2", "3"})
		require.Nil(t, errs)
	})

	t.Run("LoadManyThunk does not contain race conditions", func(t *testing.T) {
		t.Parallel()
		dl, _ := GetMockedLoader(0)
		future := dl.LoadManyThunk([]string{"1", "2", "3"})
		go future()
		go future()
	})

	t.Run("requests are batched", func(t *testing.T) {
		t.Parallel()
		dl, fetchCallsPtr := GetMockedLoader(0)
		future1 := dl.LoadThunk("1")
		future2 := dl.LoadThunk("2")

		_, err := future1()
		require.Nil(t, err)
		_, err = future2()
		require.Nil(t, err)

		fetchCalls := *fetchCallsPtr
		require.Len(t, fetchCalls, 1)
		require.Equal(t, []string{"1", "2"}, fetchCalls[0])
	})

	t.Run("number of results matches number of keys", func(t *testing.T) {
		t.Parallel()
		faultyLoader := dataloader.New(
			dataloader.Config[string, string]{Fetch: func(keys []string) (results []string, errs []error) {
				results = make([]string, len(keys)-1)
				return results, nil
			}})

		n := 10
		futures := []func() (string, error){}
		for i := 0; i < n; i++ {
			key := strconv.Itoa(i)
			futures = append(futures, faultyLoader.LoadThunk(key))
		}

		for _, future := range futures {
			_, err := future()
			require.Error(t, err)
		}
	})

	t.Run("requests are batched with max batch size", func(t *testing.T) {
		t.Parallel()
		dl, fetchCallsPtr := GetMockedLoader(2)
		future1 := dl.LoadThunk("1")
		future2 := dl.LoadThunk("2")
		future3 := dl.LoadThunk("3")

		_, err := future1()
		require.Nil(t, err)
		_, err = future2()
		require.Nil(t, err)
		_, err = future3()
		require.Nil(t, err)

		fetchCalls := *fetchCallsPtr
		require.Len(t, fetchCalls, 2)
		require.Equal(t, []string{"1", "2"}, fetchCalls[0])
		require.Equal(t, []string{"3"}, fetchCalls[1])
	})

	t.Run("repeated requests are cached", func(t *testing.T) {
		t.Parallel()
		dl, fetchCallsPtr := GetMockedLoader(0)
		future1 := dl.LoadThunk("1")
		future2 := dl.LoadThunk("1")

		_, err := future1()
		require.Nil(t, err)
		_, err = future2()
		require.Nil(t, err)

		fetchCalls := *fetchCallsPtr
		require.Len(t, fetchCalls, 1)
		require.Equal(t, []string{"1"}, fetchCalls[0])
	})

	t.Run("partial fetch", func(t *testing.T) {
		t.Parallel()
		dl, fetchCallsPtr := GetMockedLoader(0)
		{
			values, errors := dl.LoadMany([]string{"U1", "E2", "U3"})
			require.Equal(t, "U1", values[0].ID)
			require.Error(t, errors[1])
			require.Equal(t, "U3", values[2].ID)
		}
		{
			values, errors := dl.LoadMany([]string{"E2", "U3", "E4", "U5"})
			require.Error(t, errors[0])
			require.Equal(t, "U3", values[1].ID)
			require.Error(t, errors[2])
			require.Equal(t, "U5", values[3].ID)
		}

		fetchCalls := *fetchCallsPtr
		require.Len(t, fetchCalls, 2)
		require.Len(t, fetchCalls[1], 2) // Second fetch call only retrieves additional keys E4 and U5
	})
}
