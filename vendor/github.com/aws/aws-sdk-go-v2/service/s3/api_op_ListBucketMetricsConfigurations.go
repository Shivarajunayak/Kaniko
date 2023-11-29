// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Lists the metrics configurations for the bucket. The metrics configurations are
// only for the request metrics of the bucket and do not provide information on
// daily storage metrics. You can have up to 1,000 configurations per bucket. This
// action supports list pagination and does not return more than 100 configurations
// at a time. Always check the IsTruncated element in the response. If there are
// no more configurations to list, IsTruncated is set to false. If there are more
// configurations to list, IsTruncated is set to true, and there is a value in
// NextContinuationToken . You use the NextContinuationToken value to continue the
// pagination of the list by passing the value in continuation-token in the
// request to GET the next page. To use this operation, you must have permissions
// to perform the s3:GetMetricsConfiguration action. The bucket owner has this
// permission by default. The bucket owner can grant this permission to others. For
// more information about permissions, see Permissions Related to Bucket
// Subresource Operations (https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-with-s3-actions.html#using-with-s3-actions-related-to-bucket-subresources)
// and Managing Access Permissions to Your Amazon S3 Resources (https://docs.aws.amazon.com/AmazonS3/latest/userguide/s3-access-control.html)
// . For more information about metrics configurations and CloudWatch request
// metrics, see Monitoring Metrics with Amazon CloudWatch (https://docs.aws.amazon.com/AmazonS3/latest/dev/cloudwatch-monitoring.html)
// . The following operations are related to ListBucketMetricsConfigurations :
//   - PutBucketMetricsConfiguration (https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketMetricsConfiguration.html)
//   - GetBucketMetricsConfiguration (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketMetricsConfiguration.html)
//   - DeleteBucketMetricsConfiguration (https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketMetricsConfiguration.html)
func (c *Client) ListBucketMetricsConfigurations(ctx context.Context, params *ListBucketMetricsConfigurationsInput, optFns ...func(*Options)) (*ListBucketMetricsConfigurationsOutput, error) {
	if params == nil {
		params = &ListBucketMetricsConfigurationsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListBucketMetricsConfigurations", params, optFns, c.addOperationListBucketMetricsConfigurationsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListBucketMetricsConfigurationsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListBucketMetricsConfigurationsInput struct {

	// The name of the bucket containing the metrics configurations to retrieve.
	//
	// This member is required.
	Bucket *string

	// The marker that is used to continue a metrics configuration listing that has
	// been truncated. Use the NextContinuationToken from a previously truncated list
	// response to continue the listing. The continuation token is an opaque value that
	// Amazon S3 understands.
	ContinuationToken *string

	// The account ID of the expected bucket owner. If the bucket is owned by a
	// different account, the request fails with the HTTP status code 403 Forbidden
	// (access denied).
	ExpectedBucketOwner *string

	noSmithyDocumentSerde
}

func (in *ListBucketMetricsConfigurationsInput) bindEndpointParams(p *EndpointParameters) {
	p.Bucket = in.Bucket

}

type ListBucketMetricsConfigurationsOutput struct {

	// The marker that is used as a starting point for this metrics configuration list
	// response. This value is present if it was sent in the request.
	ContinuationToken *string

	// Indicates whether the returned list of metrics configurations is complete. A
	// value of true indicates that the list is not complete and the
	// NextContinuationToken will be provided for a subsequent request.
	IsTruncated *bool

	// The list of metrics configurations for a bucket.
	MetricsConfigurationList []types.MetricsConfiguration

	// The marker used to continue a metrics configuration listing that has been
	// truncated. Use the NextContinuationToken from a previously truncated list
	// response to continue the listing. The continuation token is an opaque value that
	// Amazon S3 understands.
	NextContinuationToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListBucketMetricsConfigurationsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpListBucketMetricsConfigurations{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpListBucketMetricsConfigurations{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListBucketMetricsConfigurations"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpListBucketMetricsConfigurationsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListBucketMetricsConfigurations(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addListBucketMetricsConfigurationsUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSerializeImmutableHostnameBucketMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func (v *ListBucketMetricsConfigurationsInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opListBucketMetricsConfigurations(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListBucketMetricsConfigurations",
	}
}

// getListBucketMetricsConfigurationsBucketMember returns a pointer to string
// denoting a provided bucket member valueand a boolean indicating if the input has
// a modeled bucket name,
func getListBucketMetricsConfigurationsBucketMember(input interface{}) (*string, bool) {
	in := input.(*ListBucketMetricsConfigurationsInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addListBucketMetricsConfigurationsUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getListBucketMetricsConfigurationsBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}
