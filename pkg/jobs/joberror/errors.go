// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package joberror

import (
	circuitbreaker "github.com/cockroachdb/circuitbreaker"
	"github.com/cockroachdb/cockroach/pkg/kv/kvclient/kvcoord"
	"github.com/cockroachdb/cockroach/pkg/sql/flowinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlerrors"
	"github.com/cockroachdb/cockroach/pkg/util/circuit"
	"github.com/cockroachdb/cockroach/pkg/util/grpcutil"
	"github.com/cockroachdb/cockroach/pkg/util/sysutil"
	"github.com/cockroachdb/errors"
)

// isBreakerOpenError returns true if err is a circuit.ErrBreakerOpen.
//
// NB: Two packages have ErrBreakerOpen error types.  The cicruitbreaker package
// is used by the nodedialer. The circuit package is used by kvserver.
func isBreakerOpenError(err error) bool {
	return errors.Is(err, circuit.ErrBreakerOpen) || errors.Is(err, circuitbreaker.ErrBreakerOpen)
}

// IsPermanentBulkJobError returns true if the error results in a permanent
// failure of a bulk job (IMPORT, BACKUP, RESTORE). This function is an
// allowlist instead of a blocklist: only known safe errors are confirmed to not
// be permanent errors. Anything unknown is assumed to be permanent.
func IsPermanentBulkJobError(err error) bool {
	if err == nil {
		return false
	}
	return !sqlerrors.IsDistSQLRetryableError(err) &&
		!grpcutil.IsClosedConnection(err) &&
		!flowinfra.IsFlowRetryableError(err) &&
		!flowinfra.IsNoInboundStreamConnectionError(err) &&
		!kvcoord.IsSendError(err) &&
		!isBreakerOpenError(err) &&
		!sysutil.IsErrConnectionReset(err) &&
		!sysutil.IsErrConnectionRefused(err)
}
