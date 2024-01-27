package tencent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func convertToPointerSlice(slice []string) []*string {
	// 创建一个新的切片用于存放指针
	pointerSlice := make([]*string, len(slice))

	// 遍历原始切片
	for i, value := range slice {
		// 为每个字符串创建一个指针，并将其放入新的切片中
		pointer := &value
		pointerSlice[i] = pointer
	}

	return pointerSlice
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//TODO implement me
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = common.StringPtr(tplId)
	request.TemplateParamSet = convertToPointerSlice(args)
	request.PhoneNumberSet = convertToPointerSlice(numbers)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := s.client.SendSms(request)
	// 处理异常
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("%s", b)
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr != nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			return errors.New("not send")
		}
	}
	return nil
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}
