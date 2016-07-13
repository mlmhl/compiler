package ast

import "github.com/mlmhl/compiler/common"

type Parameter struct {
	name string
	typ  Type

	location *common.Location
}
