# upload template to s3
aws s3 cp cloudformation.template s3://bt-resource/cloudformation.template
# generate a link that expires
aws s3 presign bt-resource/cloudformation.template
# update cloudformation using link to s3
aws cloudformation update-stack --stack-name bt-stack --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM --template-url https://bt-resource.s3.amazonaws.com/cloudformation.template?AWSAccessKeyId=ASIATZZT634JG7SBZ6TN&Signature=4GBnkVVpsgXvN5evbHt1ETH%2FSO0%3D&Expires=1576565225&x-amz-security-token=FwoGZXIvYXdzED8aDGhRUcz%2FSNpABImiESLwAfRsQueRBm0sk53U4aYqf%2FIDLAOdCLo%2BBc83sDnZWBKId%2FAdFbjk1gAZyvcIMWV4QC6qfR8neik2LuL7kNdIoPMhbZaSIoYt75PfZAk9hn3zUeYSdVJM1OJa3feuH8XlgUcnzZmD9wJuuR3yS58B59cMqXQ8aLmzobYZ%2BkMUu3KUJRbIeE110piv6c%2BfAaQrXSW9VwyGAswXYDY14I0eevbGCIympTZ8cOE88yWoEoSkO3p4PqpKR%2FySsSvIBKD1gvngs8Fa3NgFulb4YWXHXfxqvYqJsfx0oh6RKPaxdV735sisLKxjAPjeiL%2BlG6wcQSi9zeHvBTIrDro6%2F%2Fl5Hg35lzAJl9rAZRFFzbtm7bp1lZs%2FIQD8IDbqXRacuedXvLGcig%3D%3D
# get stack info
aws cloudformation describe-stacks --stack-name bt-stack
