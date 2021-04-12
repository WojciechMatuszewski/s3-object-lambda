# Learnings

- The Arns for the _AccessPoints_ are weird

  1. `arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point/object/Alice/`
     Notice the `/object` there - it is not a folder, it is something required by the Access Point

- **You have to use the ARN of the Object-Lambda access point, not the access point you are associating with!**

- **Before you do s3api get-object ensure that your version of CLI is the latest**. Otherwise, you might encounter errors saying that the arn is invalid.

- You **have to have `s3:ListBuckets` permissions** otherwise you will get access denied while trying to get the object through the access point.
  Sadly, this means this template will not work on ACloudGuru Sandbox :C

- I was not be able to debug the above problem with either `CloudTrial` or `CloudWatch events`. It appears that the permissions evaluation is done before the call to s3?

- How to you even front the access point?

  1. You have to provide some kind of mechanism to generate presigned URLs.
     You can reference the _AccessPoint_ via signed URL

- From what I observed, _AccessPoint_ is like a lens (with policies) that is applied to the bucket.
  The objects are not copied, it's not a different storage type or anything.

- You can upload things to the _AccessPoint_ arn, the object will still be visible in your bucket.
