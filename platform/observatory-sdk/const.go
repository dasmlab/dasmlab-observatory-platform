package observatory

// Package observatory holds shared platform constants (ADR-0001 / ADR-0002).

const SchemaVersion = "1"

// Product codes for the Observatory family.
const (
	ProductDCO   = "dco"
	ProductDNO   = "dno"
	ProductDSO   = "dso"
	ProductDPO   = "dpo"
	ProductDAO   = "dao"
	ProductDAOps = "daops"
	ProductDIO   = "dio"
	ProductDUO   = "duo"
)

// PilotTenant is the Phase-1 single tenant.
const PilotTenant = "dasmlab.org"
