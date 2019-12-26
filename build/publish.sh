ROOT="$(cd $(dirname $0)/.. && pwd)"
TARGET="${ROOT}/target"
ARTIFACT="${TARGET}/artifact"


for filename in "${ARTIFACT}"/resources/*
do
  fn=$(basename "${filename}")
  aws s3 cp "${ARTIFACT}/resources/${fn}" "s3://bt-lambdas/"$fn
done



# upload template to s3
aws s3 cp cloudformation.template s3://bt-resource/cloudformation.template
# generate a link that expires
templateurl="$(aws s3 presign bt-resource/cloudformation.template)"
# update cloudformation using link to s3
aws cloudformation update-stack --stack-name bt-stack --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM --template-url "$templateurl" --parameters ParameterKey=BigTalkBucket,ParameterValue="bt-lambdas"
# ParameterKey=PollerFunctionVersion,ParameterValue="" ParameterKey=PosterFunctionVersion,ParameterValue="" ParameterKey=QuestionsRestFunctionVersion,ParameterValue=""
# get stack info
aws cloudformation describe-stacks --stack-name bt-stack
