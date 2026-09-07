# \AuthenticationAPI

All URIs are relative to *http://my-shopping-lists.com/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**LoginPost**](AuthenticationAPI.md#LoginPost) | **Post** /login | User login



## LoginPost

> map[string]string LoginPost(ctx).MainLoginRequest(mainLoginRequest).Execute()

User login



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	mainLoginRequest := *openapiclient.NewMainLoginRequest() // MainLoginRequest | User login credentials

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthenticationAPI.LoginPost(context.Background()).MainLoginRequest(mainLoginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticationAPI.LoginPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LoginPost`: map[string]string
	fmt.Fprintf(os.Stdout, "Response from `AuthenticationAPI.LoginPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLoginPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mainLoginRequest** | [**MainLoginRequest**](MainLoginRequest.md) | User login credentials | 

### Return type

**map[string]string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

