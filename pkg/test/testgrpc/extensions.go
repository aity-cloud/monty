package testgrpc

import (
	storagev1 "github.com/aity-cloud/monty/pkg/apis/storage/v1"
)

func (t *TestSecret) RedactSecrets() {
	t.Password = storagev1.Redacted
}

func (t *TestSecret) CheckRedactedSecrets(unredacted, redacted *TestSecret) bool {
	return unredacted.Username == redacted.Username && redacted.Password == storagev1.Redacted
}
