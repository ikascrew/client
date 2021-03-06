package client

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/ikascrew/pb"
)

//const ServerIP = "172.16.10.116"
//const ServerIP = "192.168.12.1"
//const ServerIP = "192.168.11.11"

//const ServerIP = "10.0.0.1"

const ServerIP = "localhost"

func callSync() (*pb.SyncReply, error) {

	conn, err := grpc.Dial(ServerIP+":55555", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	r, err := c.Sync(context.Background(), &pb.SyncRequest{})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func callEffect(id int64, t string) error {

	conn, err := grpc.Dial(ServerIP+":55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)
	_, err = c.Effect(context.Background(), &pb.EffectRequest{
		Id:   id,
		Type: t,
	})

	if err != nil {
		return err
	}

	return nil
}

func callVolume(msg pb.VolumeMessage) error {

	conn, err := grpc.Dial(ServerIP+":55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	_, err = c.PutVolume(context.Background(), &msg)
	if err != nil {
		return err
	}
	return nil
}

func callNext() error {
	return callSwitch("next")
}

func callPrev() error {
	return callSwitch("prev")
}

func callSwitch(t string) error {

	conn, err := grpc.Dial(ServerIP+":55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	_, err = c.Switch(context.Background(), &pb.SwitchRequest{
		Type: t,
	})

	if err != nil {
		return err
	}

	return nil
}
