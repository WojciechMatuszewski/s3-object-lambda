# Learnings

- The Arns for the _AccessPoints_ are weird
  `arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point/object/Alice/`
  Notice the `/object` there - it is not a folder, it is something required by the Access Point

- Permissions: it seems like the `*` is needed for the `s3-object-lambda:WriteGetObjectResponse` action.
  Otherwise you will get 403 forbidden while calling the `WriteGetObjectResponse` API.
  The same policy is present within the documentation [here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/olap-policies.html)

  > Your AWS Lambda function needs permission to call back to the Object Lambda access point with the WriteGetObjectResponse. Add the following statement to the Execution role that is used by the Lambda function.

- If you delete the object user is requesting within the _object lambda_, the s3 will respond with xml equivalent of 404.
  What is interesting is that your _object lambda_ will be invoked either way.

- Your object lambda **must make the `WriteGetObjectResponse` call**. Otherwise you will not be able to get the object at all.

- **You have to use the ARN of the Object-Lambda access point, not the access point you are associating with!**

- **Before you do s3api get-object ensure that your version of CLI is the latest**. Otherwise, you might encounter errors saying that the arn is invalid.

- You **have to have `s3:ListBuckets` permissions** otherwise you will get access denied while trying to get the object through the access point. Sadly, this means this template will not work on ACloudGuru Sandbox :C
  I was not be able to debug the above problem with either `CloudTrial` or `CloudWatch events`. It appears that the permissions evaluation is done before the call to s3?

- _AccessPoints_ are only available to entities with some kind of identity.
  You cannot use the _object lambda_ backed _AccessPoint_ to front things for unauthenticated users (I'm not talking about Cognito here).

- From what I observed, _AccessPoint_ is like a lens (with policies) that is applied to the bucket.
  The objects are not copied, it's not a different storage type or anything.

- You can upload things to the _AccessPoint_ arn, the object will still be visible in your bucket.
