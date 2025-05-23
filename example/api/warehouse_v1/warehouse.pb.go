// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.3
// source: proto/warehouse/v1/warehouse.proto

package warehouse_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// *
// ProductInventoryは、扱っている商品の情報を表します。
type ProductInventory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 商品番号
	Number string `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	// 商品名
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// 価格（円）
	Price uint32 `protobuf:"varint,3,opt,name=price,proto3" json:"price,omitempty"`
	// 在庫数
	QuantityAvailable uint32 `protobuf:"varint,4,opt,name=quantity_available,json=quantityAvailable,proto3" json:"quantity_available,omitempty"`
}

func (x *ProductInventory) Reset() {
	*x = ProductInventory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductInventory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductInventory) ProtoMessage() {}

func (x *ProductInventory) ProtoReflect() protoreflect.Message {
	mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductInventory.ProtoReflect.Descriptor instead.
func (*ProductInventory) Descriptor() ([]byte, []int) {
	return file_proto_warehouse_v1_warehouse_proto_rawDescGZIP(), []int{0}
}

func (x *ProductInventory) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *ProductInventory) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProductInventory) GetPrice() uint32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ProductInventory) GetQuantityAvailable() uint32 {
	if x != nil {
		return x.QuantityAvailable
	}
	return 0
}

// *
// ListProductInventoriesは、扱っている商品一覧を返します。
// - num_of_productsは、一覧として返す最大個数を指定します。
// -
// page_tokenには、空文字列もしくはレスポンスで返されるnext_page_tokenを設定します。
//   - 空文字列の場合は、商品一覧の最初から返されます。
//
// - 商品一覧は、nameのアルファベット順に昇順に返されます。
//
// [エラー]
// - InvalidArgument:
//   - num_of_productsが0
//   - page_tokenが不正な値
type ListProductInventoriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumOfProducts uint32 `protobuf:"varint,1,opt,name=num_of_products,json=numOfProducts,proto3" json:"num_of_products,omitempty"`
	PageToken     string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListProductInventoriesRequest) Reset() {
	*x = ListProductInventoriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProductInventoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductInventoriesRequest) ProtoMessage() {}

func (x *ListProductInventoriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductInventoriesRequest.ProtoReflect.Descriptor instead.
func (*ListProductInventoriesRequest) Descriptor() ([]byte, []int) {
	return file_proto_warehouse_v1_warehouse_proto_rawDescGZIP(), []int{1}
}

func (x *ListProductInventoriesRequest) GetNumOfProducts() uint32 {
	if x != nil {
		return x.NumOfProducts
	}
	return 0
}

func (x *ListProductInventoriesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListProductInventoriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductInventories []*ProductInventory `protobuf:"bytes,1,rep,name=product_inventories,json=productInventories,proto3" json:"product_inventories,omitempty"`
	// 商品がさらに存在しない場合、空文字列が返されます。
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListProductInventoriesResponse) Reset() {
	*x = ListProductInventoriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProductInventoriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductInventoriesResponse) ProtoMessage() {}

func (x *ListProductInventoriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductInventoriesResponse.ProtoReflect.Descriptor instead.
func (*ListProductInventoriesResponse) Descriptor() ([]byte, []int) {
	return file_proto_warehouse_v1_warehouse_proto_rawDescGZIP(), []int{2}
}

func (x *ListProductInventoriesResponse) GetProductInventories() []*ProductInventory {
	if x != nil {
		return x.ProductInventories
	}
	return nil
}

func (x *ListProductInventoriesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// *
// ShipProductは、指定された商品を出荷します。
//   - Shippingサービスでステータスを問い合わせるための
//     shipping_idを返します。
//
// [エラー]
// - InvalidArgument:
//   - order_idが空文字列
//   - numberが空文字列
//   - num_of_itemsが0
//   - shipping_addressが空文字列
//
// - NotFound:
//   - numberで指定された商品は扱っていない
//
// - FailedPrecondition:
//   - num_of_itemsで指定された個数の在庫がないため出荷できない。
type ShipProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 注文番号
	OrderId string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	// 商品番号
	Number string `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	// 個数
	NumOfItems uint32 `protobuf:"varint,3,opt,name=num_of_items,json=numOfItems,proto3" json:"num_of_items,omitempty"`
	// 送付先
	ShippingAddress string `protobuf:"bytes,4,opt,name=shipping_address,json=shippingAddress,proto3" json:"shipping_address,omitempty"`
}

func (x *ShipProductRequest) Reset() {
	*x = ShipProductRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShipProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShipProductRequest) ProtoMessage() {}

func (x *ShipProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShipProductRequest.ProtoReflect.Descriptor instead.
func (*ShipProductRequest) Descriptor() ([]byte, []int) {
	return file_proto_warehouse_v1_warehouse_proto_rawDescGZIP(), []int{3}
}

func (x *ShipProductRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *ShipProductRequest) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *ShipProductRequest) GetNumOfItems() uint32 {
	if x != nil {
		return x.NumOfItems
	}
	return 0
}

func (x *ShipProductRequest) GetShippingAddress() string {
	if x != nil {
		return x.ShippingAddress
	}
	return ""
}

type ShipProductResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 出荷番号
	ShippingId string `protobuf:"bytes,1,opt,name=shipping_id,json=shippingId,proto3" json:"shipping_id,omitempty"`
}

func (x *ShipProductResponse) Reset() {
	*x = ShipProductResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShipProductResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShipProductResponse) ProtoMessage() {}

func (x *ShipProductResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_warehouse_v1_warehouse_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShipProductResponse.ProtoReflect.Descriptor instead.
func (*ShipProductResponse) Descriptor() ([]byte, []int) {
	return file_proto_warehouse_v1_warehouse_proto_rawDescGZIP(), []int{4}
}

func (x *ShipProductResponse) GetShippingId() string {
	if x != nil {
		return x.ShippingId
	}
	return ""
}

var File_proto_warehouse_v1_warehouse_proto protoreflect.FileDescriptor

var file_proto_warehouse_v1_warehouse_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2f, 0x79, 0x6f, 0x73, 0x68, 0x69, 0x6b, 0x69, 0x73, 0x68, 0x69,
	0x62, 0x61, 0x74, 0x61, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75,
	0x73, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x83, 0x01, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x2d, 0x0a, 0x12,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x66, 0x0a, 0x1d, 0x4c,
	0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0f,
	0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0xbc, 0x01, 0x0a, 0x1e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x72, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x5f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x41, 0x2e, 0x79, 0x6f, 0x73, 0x68, 0x69, 0x6b, 0x69, 0x73, 0x68, 0x69,
	0x62, 0x61, 0x74, 0x61, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75,
	0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x94, 0x01, 0x0a, 0x12, 0x53, 0x68, 0x69, 0x70, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0c,
	0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x29,
	0x0a, 0x10, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69,
	0x6e, 0x67, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x36, 0x0a, 0x13, 0x53, 0x68, 0x69,
	0x70, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x49,
	0x64, 0x32, 0xe6, 0x02, 0x0a, 0x09, 0x57, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x12,
	0xbb, 0x01, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x4e, 0x2e, 0x79, 0x6f, 0x73,
	0x68, 0x69, 0x6b, 0x69, 0x73, 0x68, 0x69, 0x62, 0x61, 0x74, 0x61, 0x2e, 0x63, 0x6f, 0x75, 0x72,
	0x69, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x4f, 0x2e, 0x79, 0x6f, 0x73,
	0x68, 0x69, 0x6b, 0x69, 0x73, 0x68, 0x69, 0x62, 0x61, 0x74, 0x61, 0x2e, 0x63, 0x6f, 0x75, 0x72,
	0x69, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x9a, 0x01,
	0x0a, 0x0b, 0x53, 0x68, 0x69, 0x70, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x43, 0x2e,
	0x79, 0x6f, 0x73, 0x68, 0x69, 0x6b, 0x69, 0x73, 0x68, 0x69, 0x62, 0x61, 0x74, 0x61, 0x2e, 0x63,
	0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x68, 0x69, 0x70, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x44, 0x2e, 0x79, 0x6f, 0x73, 0x68, 0x69, 0x6b, 0x69, 0x73, 0x68, 0x69, 0x62,
	0x61, 0x74, 0x61, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x69, 0x70, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2f, 0x77,
	0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_warehouse_v1_warehouse_proto_rawDescOnce sync.Once
	file_proto_warehouse_v1_warehouse_proto_rawDescData = file_proto_warehouse_v1_warehouse_proto_rawDesc
)

func file_proto_warehouse_v1_warehouse_proto_rawDescGZIP() []byte {
	file_proto_warehouse_v1_warehouse_proto_rawDescOnce.Do(func() {
		file_proto_warehouse_v1_warehouse_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_warehouse_v1_warehouse_proto_rawDescData)
	})
	return file_proto_warehouse_v1_warehouse_proto_rawDescData
}

var file_proto_warehouse_v1_warehouse_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_warehouse_v1_warehouse_proto_goTypes = []interface{}{
	(*ProductInventory)(nil),               // 0: yoshikishibata.courier.example.api.warehouse.v1.ProductInventory
	(*ListProductInventoriesRequest)(nil),  // 1: yoshikishibata.courier.example.api.warehouse.v1.ListProductInventoriesRequest
	(*ListProductInventoriesResponse)(nil), // 2: yoshikishibata.courier.example.api.warehouse.v1.ListProductInventoriesResponse
	(*ShipProductRequest)(nil),             // 3: yoshikishibata.courier.example.api.warehouse.v1.ShipProductRequest
	(*ShipProductResponse)(nil),            // 4: yoshikishibata.courier.example.api.warehouse.v1.ShipProductResponse
}
var file_proto_warehouse_v1_warehouse_proto_depIdxs = []int32{
	0, // 0: yoshikishibata.courier.example.api.warehouse.v1.ListProductInventoriesResponse.product_inventories:type_name -> yoshikishibata.courier.example.api.warehouse.v1.ProductInventory
	1, // 1: yoshikishibata.courier.example.api.warehouse.v1.Warehouse.ListProductInventories:input_type -> yoshikishibata.courier.example.api.warehouse.v1.ListProductInventoriesRequest
	3, // 2: yoshikishibata.courier.example.api.warehouse.v1.Warehouse.ShipProduct:input_type -> yoshikishibata.courier.example.api.warehouse.v1.ShipProductRequest
	2, // 3: yoshikishibata.courier.example.api.warehouse.v1.Warehouse.ListProductInventories:output_type -> yoshikishibata.courier.example.api.warehouse.v1.ListProductInventoriesResponse
	4, // 4: yoshikishibata.courier.example.api.warehouse.v1.Warehouse.ShipProduct:output_type -> yoshikishibata.courier.example.api.warehouse.v1.ShipProductResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_warehouse_v1_warehouse_proto_init() }
func file_proto_warehouse_v1_warehouse_proto_init() {
	if File_proto_warehouse_v1_warehouse_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_warehouse_v1_warehouse_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProductInventory); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_warehouse_v1_warehouse_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListProductInventoriesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_warehouse_v1_warehouse_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListProductInventoriesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_warehouse_v1_warehouse_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShipProductRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_warehouse_v1_warehouse_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShipProductResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_warehouse_v1_warehouse_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_warehouse_v1_warehouse_proto_goTypes,
		DependencyIndexes: file_proto_warehouse_v1_warehouse_proto_depIdxs,
		MessageInfos:      file_proto_warehouse_v1_warehouse_proto_msgTypes,
	}.Build()
	File_proto_warehouse_v1_warehouse_proto = out.File
	file_proto_warehouse_v1_warehouse_proto_rawDesc = nil
	file_proto_warehouse_v1_warehouse_proto_goTypes = nil
	file_proto_warehouse_v1_warehouse_proto_depIdxs = nil
}
