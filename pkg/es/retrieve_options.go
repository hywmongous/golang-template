package es

import (
	"errors"

	merr "github.com/hywmongous/example-service/pkg/errors"
)

type RetrieveOptions struct {
	MinVersion   *EventVersion
	MaxVersion   *EventVersion
	MinTimestamp *EventTimestamp
	MaxTimestamp *EventTimestamp
}

var (
	ErrOptionIsNil = errors.New("option is nil")

	ErrMinVersionGreaterThanMax   = errors.New("min version cannot be greater than the max")
	ErrMaxVersionLessThanMin      = errors.New("max version cannot be less than max")
	ErrMinTimestampGreaterThanMax = errors.New("min timestamp cannot be greater than the max")
	ErrMaxTimestampLessThanMin    = errors.New("max timestamp cannot be less than max")
)

func CreateRetrieveOptions() RetrieveOptions {
	return RetrieveOptions{}
}

func (options *RetrieveOptions) SetMinVersion(version EventVersion) error {
	if version > *options.MaxVersion {
		return merr.CreateFailedStructInvocation("RetrieveOptions", "SetMinVersion", ErrMinTimestampGreaterThanMax)
	}
	options.MinVersion = &version
	return nil
}

func (options *RetrieveOptions) SetMaxVersion(version EventVersion) error {
	if version < *options.MinVersion {
		return merr.CreateFailedStructInvocation("RetrieveOptions", "SetMinVersion", ErrMaxVersionLessThanMin)
	}
	options.MaxVersion = &version
	return nil
}

func (options *RetrieveOptions) SetMinTimestamp(timestamp EventTimestamp) error {
	if timestamp > *options.MaxTimestamp {
		return merr.CreateFailedStructInvocation("RetrieveOptions", "SetMinVersion", ErrMinTimestampGreaterThanMax)
	}
	options.MinTimestamp = &timestamp
	return nil
}

func (options *RetrieveOptions) SetMaxTimestamp(timestamp EventTimestamp) error {
	if timestamp < *options.MinTimestamp {
		return merr.CreateFailedStructInvocation("RetrieveOptions", "SetMinVersion", ErrMaxVersionLessThanMin)
	}
	options.MaxTimestamp = &timestamp
	return nil
}

func MergeRetrieveOptions(opts ...*RetrieveOptions) (RetrieveOptions, error) {
	options := CreateRetrieveOptions()
	for _, opt := range opts {
		if opt == nil {
			return CreateRetrieveOptions(), ErrOptionIsNil
		}

		if opt.MinVersion != nil {
			if err := options.SetMinVersion(*opt.MinVersion); err != nil {
				return CreateRetrieveOptions(), merr.CreateFailedInvocation("MergeRetrieveOptions", err)
			}
		}

		if opt.MaxVersion != nil {
			if err := options.SetMaxVersion(*opt.MaxVersion); err != nil {
				return CreateRetrieveOptions(), merr.CreateFailedInvocation("MergeRetrieveOptions", err)
			}
		}

		if opt.MinTimestamp != nil {
			if err := options.SetMinTimestamp(*opt.MinTimestamp); err != nil {
				return CreateRetrieveOptions(), merr.CreateFailedInvocation("MergeRetrieveOptions", err)
			}
		}

		if opt.MaxTimestamp != nil {
			if err := options.SetMaxTimestamp(*opt.MaxTimestamp); err != nil {
				return CreateRetrieveOptions(), merr.CreateFailedInvocation("MergeRetrieveOptions", err)
			}
		}
	}
	return options, nil
}
