package loadbalance_test

import (
	"testing"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"

	"github.com/stretchr/testify/require"

	"github.com/qs-lzh/proglog/internal/loadbalance"
)

func TestPickerNoSubConnAvailable(t *testing.T) {
	picker := &loadbalance.Picker{}
	for _, method := range []string{"/log.vX.Log/Produce", "/log.vX.Log/Consume"} {
		info := balancer.PickInfo{FullMethodName: method}
		result, err := picker.Pick(info)
		require.Equal(t, balancer.ErrNoSubConnAvailable, err)
		require.Nil(t, result.SubConn)
	}
}

func TestPickerProducesToLeader(t *testing.T) {
	picker, subConns := setupTest()
	info := balancer.PickInfo{FullMethodName: "/log.vX.Log/Produce"}
	for range 5 {
		gotPick, err := picker.Pick(info)
		require.NoError(t, err)
		require.Equal(t, subConns[0], gotPick.SubConn)
	}
}

func TestPickerConsumesFromFollowers(t *testing.T) {
	picker, subConns := setupTest()
	info := balancer.PickInfo{
		FullMethodName: "/log.vX.Log/Consume",
	}
	for i := range 5 {
		pick, err := picker.Pick(info)
		require.NoError(t, err)
		require.Equal(t, subConns[i%2+1], pick.SubConn)
	}
}

func setupTest() (*loadbalance.Picker, []*subConn) {
	var subConns []*subConn
	buildInfo := base.PickerBuildInfo{ReadySCs: make(map[balancer.SubConn]base.SubConnInfo)}
	for i := range 3 {
		sc := &subConn{}
		addr := resolver.Address{Attributes: attributes.New("is_leader", i == 0)}
		// NOTE the key type of balancer.base.PickerBuildInfo.ReadySCs -- balancer.SubConn interface has changed.
		// Now it has an internal interface, which means we can not satisfy  balancer.SubConn by ourselves
		// So the code in the book does not work anymore, and I don't know how to fix it
		// And I do not change the code at page 192, 193 of the book
		sc.UpdateAddresses([]resolver.Address{addr})
		buildInfo.ReadySCs[sc] = base.SubConnInfo{Address: addr}
		subConns = append(subConns, sc)
	}
	picker := &loadbalance.Picker{}
	picker.Build(buildInfo)
	return picker, subConns
}

type subConn struct {
	addrs []resolver.Address
}

func (s *subConn) UpdateAddresses(addrs []resolver.Address) {
	s.addrs = addrs
}

func (s *subConn) Connect() {}
