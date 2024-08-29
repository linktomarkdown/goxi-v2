package goxi_v2

import (
	"context"
	"fmt"
	"github.com/minio/madmin-go/v3"
	"log"
)

type MDM struct {
	*madmin.AdminClient
	ctx context.Context
}

func NewMDM(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MDM, error) {
	log.Printf("初始化连接Minio Admin" + endpoint + " " + accessKeyID + " " + secretAccessKey + " " + fmt.Sprintf("%t", useSSL))
	ctx := context.Background()
	mdm, err := madmin.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Fetch service status.
	st, err := mdm.ServerInfo(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(st)
	log.Printf("初始化连接Minio Admin成功")
	return &MDM{
		AdminClient: mdm,
		ctx:         ctx,
	}, nil
}

// CreateUser 创建用户
func (m *MDM) CreateUser(accessKey, secretKey string) error {
	err := m.AddUser(m.ctx, accessKey, secretKey)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully added user %s\n", accessKey)
	return nil
}

// SetUserPolicy 设置策略
func (m *MDM) SetUserPolicy(accessKey, policyName string) error {
	err := m.SetPolicy(m.ctx, policyName, accessKey, false)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully set policy %s\n", policyName)
	return nil
}

// UserInfo 获取用户信息
func (m *MDM) UserInfo(accessKey string) (madmin.UserInfo, error) {
	userInfo, err := m.GetUserInfo(m.ctx, accessKey)
	if err != nil {
		log.Fatal(err.Error())
		return userInfo, err
	}
	return userInfo, nil
}

// PoliciesList 获取所有策略
func (m *MDM) PoliciesList() ([]madmin.KMSPolicyInfo, error) {
	policies, err := m.ListPolicies(m.ctx, "readwrite")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return policies, nil
}
