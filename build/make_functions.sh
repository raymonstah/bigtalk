ROOT="$(cd $(dirname $0)/.. && pwd)"
TARGET="${ROOT}/target"
ARTIFACT="${TARGET}/artifact"

if [[ -d ${ROOT}/functions ]] ; then
  mkdir -p "${ARTIFACT}/resources"
  for dir in "${ROOT}"/functions/*
  do
    fn=$(basename "${dir}")
    (cd "${TARGET}"; GOOS=linux go build -o "${fn}" "${ROOT}/functions/${fn}"/*.go)
    (cd "${TARGET}"; zip -r "${ARTIFACT}/resources/${fn}.zip" "${fn}")
  done
fi