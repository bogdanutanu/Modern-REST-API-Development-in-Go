# \ListsAPI

All URIs are relative to *http://my-shopping-lists.com/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListsGet**](ListsAPI.md#ListsGet) | **Get** /lists | List all shopping lists
[**ListsIdDelete**](ListsAPI.md#ListsIdDelete) | **Delete** /lists/{id} | Delete a shopping list
[**ListsIdGet**](ListsAPI.md#ListsIdGet) | **Get** /lists/{id} | Get a shopping list by ID
[**ListsIdPatch**](ListsAPI.md#ListsIdPatch) | **Patch** /lists/{id} | Partially update a shopping list
[**ListsIdPushPost**](ListsAPI.md#ListsIdPushPost) | **Post** /lists/{id}/push | Add an item to a shopping list
[**ListsIdPut**](ListsAPI.md#ListsIdPut) | **Put** /lists/{id} | Update a shopping list
[**ListsPost**](ListsAPI.md#ListsPost) | **Post** /lists | Create a new shopping list



## ListsGet

> []MainShoppingList ListsGet(ctx).Execute()

List all shopping lists



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ListsAPI.ListsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListsGet`: []MainShoppingList
	fmt.Fprintf(os.Stdout, "Response from `ListsAPI.ListsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListsGetRequest struct via the builder pattern


### Return type

[**[]MainShoppingList**](MainShoppingList.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListsIdDelete

> ListsIdDelete(ctx, id).Execute()

Delete a shopping list



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
	id := "id_example" // string | Shopping list ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ListsAPI.ListsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Shopping list ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListsIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListsIdGet

> MainShoppingList ListsIdGet(ctx, id).Execute()

Get a shopping list by ID



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
	id := "id_example" // string | Shopping list ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ListsAPI.ListsIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListsIdGet`: MainShoppingList
	fmt.Fprintf(os.Stdout, "Response from `ListsAPI.ListsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Shopping list ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**MainShoppingList**](MainShoppingList.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListsIdPatch

> MainShoppingList ListsIdPatch(ctx, id).MainShoppingListPatch(mainShoppingListPatch).Execute()

Partially update a shopping list



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
	id := "id_example" // string | Shopping list ID
	mainShoppingListPatch := *openapiclient.NewMainShoppingListPatch() // MainShoppingListPatch | Partial shopping list data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ListsAPI.ListsIdPatch(context.Background(), id).MainShoppingListPatch(mainShoppingListPatch).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListsIdPatch`: MainShoppingList
	fmt.Fprintf(os.Stdout, "Response from `ListsAPI.ListsIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Shopping list ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListsIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **mainShoppingListPatch** | [**MainShoppingListPatch**](MainShoppingListPatch.md) | Partial shopping list data | 

### Return type

[**MainShoppingList**](MainShoppingList.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListsIdPushPost

> MainShoppingList ListsIdPushPost(ctx, id).MainListPushAction(mainListPushAction).Execute()

Add an item to a shopping list



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
	id := "id_example" // string | Shopping list ID
	mainListPushAction := *openapiclient.NewMainListPushAction() // MainListPushAction | Item to add to the list

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ListsAPI.ListsIdPushPost(context.Background(), id).MainListPushAction(mainListPushAction).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsIdPushPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListsIdPushPost`: MainShoppingList
	fmt.Fprintf(os.Stdout, "Response from `ListsAPI.ListsIdPushPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Shopping list ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListsIdPushPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **mainListPushAction** | [**MainListPushAction**](MainListPushAction.md) | Item to add to the list | 

### Return type

[**MainShoppingList**](MainShoppingList.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListsIdPut

> MainShoppingList ListsIdPut(ctx, id).MainShoppingList(mainShoppingList).Execute()

Update a shopping list



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
	id := "id_example" // string | Shopping list ID
	mainShoppingList := *openapiclient.NewMainShoppingList() // MainShoppingList | Updated shopping list data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ListsAPI.ListsIdPut(context.Background(), id).MainShoppingList(mainShoppingList).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListsIdPut`: MainShoppingList
	fmt.Fprintf(os.Stdout, "Response from `ListsAPI.ListsIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Shopping list ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListsIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **mainShoppingList** | [**MainShoppingList**](MainShoppingList.md) | Updated shopping list data | 

### Return type

[**MainShoppingList**](MainShoppingList.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListsPost

> MainShoppingList ListsPost(ctx).MainShoppingList(mainShoppingList).Execute()

Create a new shopping list



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
	mainShoppingList := *openapiclient.NewMainShoppingList() // MainShoppingList | Shopping list data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ListsAPI.ListsPost(context.Background()).MainShoppingList(mainShoppingList).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ListsAPI.ListsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListsPost`: MainShoppingList
	fmt.Fprintf(os.Stdout, "Response from `ListsAPI.ListsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mainShoppingList** | [**MainShoppingList**](MainShoppingList.md) | Shopping list data | 

### Return type

[**MainShoppingList**](MainShoppingList.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

