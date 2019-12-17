ROOT="$(cd $(dirname $0)/.. && pwd)"
TARGET="${ROOT}/target"
ARTIFACT="${TARGET}/artifact"


for filename in "${ARTIFACT}"/resources/*
do
  fn=$(basename "${filename}")
  aws s3 cp "${ARTIFACT}/resources/${fn}" "s3://bt-lambdas/"$fn
done
