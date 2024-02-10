package beacon

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

func GetHistoricalSummaries(state *spec.VersionedBeaconState) ([]*capella.HistoricalSummary, error) {
	switch state.Version {
	case spec.DataVersionCapella:
		return state.Capella.HistoricalSummaries, nil
	case spec.DataVersionDeneb:
		return state.Deneb.HistoricalSummaries, nil
	default:
		return nil, errors.New("unsupported beacon state version")
	}
}

func GetGenesisTime(state *spec.VersionedBeaconState) (uint64, error) {
	switch state.Version {
	case spec.DataVersionCapella:
		return state.Capella.GenesisTime, nil
	case spec.DataVersionDeneb:
		return state.Deneb.GenesisTime, nil
	default:
		return 0, errors.New("unsupported beacon state version")
	}
}

func GetBlockRoots(beaconState spec.VersionedBeaconState) []phase0.Root {
	blockRoots := make([]phase0.Root, 0)
	switch beaconState.Version {
	case spec.DataVersionDeneb:
		blockRoots = beaconState.Deneb.BlockRoots
	case spec.DataVersionCapella:
		blockRoots = beaconState.Capella.BlockRoots
	}
	return blockRoots
}
func CreateVersionedSignedBlock(block interface{}) (spec.VersionedSignedBeaconBlock, error) {
	var versionedBlock spec.VersionedSignedBeaconBlock

	switch s := block.(type) {
	case deneb.BeaconBlock:
		var signedBlock deneb.SignedBeaconBlock
		signedBlock.Message = &s
		versionedBlock.Deneb = &signedBlock
		versionedBlock.Version = spec.DataVersionDeneb
	case capella.BeaconBlock:
		var signedBlock capella.SignedBeaconBlock
		signedBlock.Message = &s
		versionedBlock.Capella = &signedBlock
		versionedBlock.Version = spec.DataVersionCapella
	default:
		return versionedBlock, errors.New("unsupported beacon block version")
	}
	return versionedBlock, nil
}

func CreateVersionedState(state interface{}) (spec.VersionedBeaconState, error) {
	var versionedState spec.VersionedBeaconState

	switch s := state.(type) {
	case *deneb.BeaconState:
		versionedState.Deneb = s
		versionedState.Version = spec.DataVersionDeneb
	case *capella.BeaconState:
		versionedState.Capella = s
		versionedState.Version = spec.DataVersionCapella
	default:
		return versionedState, errors.New("unsupported beacon state version")
	}
	return versionedState, nil
}

func UnmarshalSSZVersionedBeaconState(data []byte) (*spec.VersionedBeaconState, error) {
	beaconState := &spec.VersionedBeaconState{}
	// Try to unmarshal using Deneb
	err := beaconState.Deneb.UnmarshalSSZ(data)
	if err != nil {
		// If Deneb fails, try Capella
		err = beaconState.Capella.UnmarshalSSZ(data)
		if err != nil {
			return nil, err
		}
	}

	return beaconState, nil
}

func MarshalSSZVersionedBeaconState(beaconState spec.VersionedBeaconState) ([]byte, error) {
	var data []byte
	// Try to marshal using Deneb
	data, err := beaconState.Deneb.MarshalSSZ()
	if err != nil {
		// If Deneb fails, try Capella
		data, err = beaconState.Capella.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
