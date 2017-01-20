package ast

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

//
// BUGS: the escaped forward solidus (`\/`) is not currently handled for strings.
//

// currentLocation converts the parser context to a Location object.
func currentLocation(c *current) *Location {
	// TODO(tsandall): is it possible to access the filename from inside the parser?
	return NewLocation(c.text, "", c.pos.line, c.pos.col)
}

func ifaceSliceToByteSlice(i interface{}) []byte {
	var buf bytes.Buffer
	for _, x := range i.([]interface{}) {
		buf.Write(x.([]byte))
	}
	return buf.Bytes()
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Program",
			pos:  position{line: 24, col: 1, offset: 547},
			expr: &actionExpr{
				pos: position{line: 24, col: 12, offset: 558},
				run: (*parser).callonProgram1,
				expr: &seqExpr{
					pos: position{line: 24, col: 12, offset: 558},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 24, col: 12, offset: 558},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 24, col: 14, offset: 560},
							label: "vals",
							expr: &zeroOrOneExpr{
								pos: position{line: 24, col: 19, offset: 565},
								expr: &seqExpr{
									pos: position{line: 24, col: 20, offset: 566},
									exprs: []interface{}{
										&labeledExpr{
											pos:   position{line: 24, col: 20, offset: 566},
											label: "head",
											expr: &ruleRefExpr{
												pos:  position{line: 24, col: 25, offset: 571},
												name: "Stmt",
											},
										},
										&labeledExpr{
											pos:   position{line: 24, col: 30, offset: 576},
											label: "tail",
											expr: &zeroOrMoreExpr{
												pos: position{line: 24, col: 35, offset: 581},
												expr: &seqExpr{
													pos: position{line: 24, col: 36, offset: 582},
													exprs: []interface{}{
														&choiceExpr{
															pos: position{line: 24, col: 37, offset: 583},
															alternatives: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 24, col: 37, offset: 583},
																	name: "ws",
																},
																&ruleRefExpr{
																	pos:  position{line: 24, col: 42, offset: 588},
																	name: "ParseError",
																},
															},
														},
														&ruleRefExpr{
															pos:  position{line: 24, col: 54, offset: 600},
															name: "Stmt",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 63, offset: 609},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 24, col: 65, offset: 611},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Stmt",
			pos:  position{line: 42, col: 1, offset: 948},
			expr: &actionExpr{
				pos: position{line: 42, col: 9, offset: 956},
				run: (*parser).callonStmt1,
				expr: &labeledExpr{
					pos:   position{line: 42, col: 9, offset: 956},
					label: "val",
					expr: &choiceExpr{
						pos: position{line: 42, col: 14, offset: 961},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 42, col: 14, offset: 961},
								name: "Package",
							},
							&ruleRefExpr{
								pos:  position{line: 42, col: 24, offset: 971},
								name: "Import",
							},
							&ruleRefExpr{
								pos:  position{line: 42, col: 33, offset: 980},
								name: "Rule",
							},
							&ruleRefExpr{
								pos:  position{line: 42, col: 40, offset: 987},
								name: "Body",
							},
							&ruleRefExpr{
								pos:  position{line: 42, col: 47, offset: 994},
								name: "Comment",
							},
							&ruleRefExpr{
								pos:  position{line: 42, col: 57, offset: 1004},
								name: "ParseError",
							},
						},
					},
				},
			},
		},
		{
			name: "ParseError",
			pos:  position{line: 51, col: 1, offset: 1368},
			expr: &actionExpr{
				pos: position{line: 51, col: 15, offset: 1382},
				run: (*parser).callonParseError1,
				expr: &anyMatcher{
					line: 51, col: 15, offset: 1382,
				},
			},
		},
		{
			name: "Package",
			pos:  position{line: 55, col: 1, offset: 1455},
			expr: &actionExpr{
				pos: position{line: 55, col: 12, offset: 1466},
				run: (*parser).callonPackage1,
				expr: &seqExpr{
					pos: position{line: 55, col: 12, offset: 1466},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 55, col: 12, offset: 1466},
							val:        "package",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 22, offset: 1476},
							name: "ws",
						},
						&labeledExpr{
							pos:   position{line: 55, col: 25, offset: 1479},
							label: "val",
							expr: &choiceExpr{
								pos: position{line: 55, col: 30, offset: 1484},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 55, col: 30, offset: 1484},
										name: "Ref",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 36, offset: 1490},
										name: "Var",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Import",
			pos:  position{line: 91, col: 1, offset: 2871},
			expr: &actionExpr{
				pos: position{line: 91, col: 11, offset: 2881},
				run: (*parser).callonImport1,
				expr: &seqExpr{
					pos: position{line: 91, col: 11, offset: 2881},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 91, col: 11, offset: 2881},
							val:        "import",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 91, col: 20, offset: 2890},
							name: "ws",
						},
						&labeledExpr{
							pos:   position{line: 91, col: 23, offset: 2893},
							label: "path",
							expr: &choiceExpr{
								pos: position{line: 91, col: 29, offset: 2899},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 91, col: 29, offset: 2899},
										name: "Ref",
									},
									&ruleRefExpr{
										pos:  position{line: 91, col: 35, offset: 2905},
										name: "Var",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 91, col: 40, offset: 2910},
							label: "alias",
							expr: &zeroOrOneExpr{
								pos: position{line: 91, col: 46, offset: 2916},
								expr: &seqExpr{
									pos: position{line: 91, col: 47, offset: 2917},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 91, col: 47, offset: 2917},
											name: "ws",
										},
										&litMatcher{
											pos:        position{line: 91, col: 50, offset: 2920},
											val:        "as",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 91, col: 55, offset: 2925},
											name: "ws",
										},
										&ruleRefExpr{
											pos:  position{line: 91, col: 58, offset: 2928},
											name: "Var",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 107, col: 1, offset: 3378},
			expr: &actionExpr{
				pos: position{line: 107, col: 9, offset: 3386},
				run: (*parser).callonRule1,
				expr: &seqExpr{
					pos: position{line: 107, col: 9, offset: 3386},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 107, col: 9, offset: 3386},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 107, col: 14, offset: 3391},
								name: "Var",
							},
						},
						&labeledExpr{
							pos:   position{line: 107, col: 18, offset: 3395},
							label: "key",
							expr: &zeroOrOneExpr{
								pos: position{line: 107, col: 22, offset: 3399},
								expr: &seqExpr{
									pos: position{line: 107, col: 24, offset: 3401},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 107, col: 24, offset: 3401},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 107, col: 26, offset: 3403},
											val:        "[",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 30, offset: 3407},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 32, offset: 3409},
											name: "Term",
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 37, offset: 3414},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 107, col: 39, offset: 3416},
											val:        "]",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 43, offset: 3420},
											name: "_",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 107, col: 48, offset: 3425},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 107, col: 54, offset: 3431},
								expr: &seqExpr{
									pos: position{line: 107, col: 56, offset: 3433},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 107, col: 56, offset: 3433},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 107, col: 58, offset: 3435},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 62, offset: 3439},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 64, offset: 3441},
											name: "Term",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 107, col: 72, offset: 3449},
							label: "body",
							expr: &seqExpr{
								pos: position{line: 107, col: 79, offset: 3456},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 107, col: 79, offset: 3456},
										name: "_",
									},
									&litMatcher{
										pos:        position{line: 107, col: 81, offset: 3458},
										val:        ":-",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 107, col: 86, offset: 3463},
										name: "_",
									},
									&ruleRefExpr{
										pos:  position{line: 107, col: 88, offset: 3465},
										name: "Body",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Body",
			pos:  position{line: 163, col: 1, offset: 5096},
			expr: &actionExpr{
				pos: position{line: 163, col: 9, offset: 5104},
				run: (*parser).callonBody1,
				expr: &seqExpr{
					pos: position{line: 163, col: 9, offset: 5104},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 163, col: 9, offset: 5104},
							label: "head",
							expr: &ruleRefExpr{
								pos:  position{line: 163, col: 14, offset: 5109},
								name: "Expr",
							},
						},
						&labeledExpr{
							pos:   position{line: 163, col: 19, offset: 5114},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 163, col: 24, offset: 5119},
								expr: &seqExpr{
									pos: position{line: 163, col: 26, offset: 5121},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 163, col: 26, offset: 5121},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 163, col: 28, offset: 5123},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 163, col: 32, offset: 5127},
											name: "_",
										},
										&choiceExpr{
											pos: position{line: 163, col: 35, offset: 5130},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 163, col: 35, offset: 5130},
													name: "Expr",
												},
												&ruleRefExpr{
													pos:  position{line: 163, col: 42, offset: 5137},
													name: "ParseError",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 173, col: 1, offset: 5357},
			expr: &actionExpr{
				pos: position{line: 173, col: 9, offset: 5365},
				run: (*parser).callonExpr1,
				expr: &seqExpr{
					pos: position{line: 173, col: 9, offset: 5365},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 173, col: 9, offset: 5365},
							label: "neg",
							expr: &zeroOrOneExpr{
								pos: position{line: 173, col: 13, offset: 5369},
								expr: &seqExpr{
									pos: position{line: 173, col: 15, offset: 5371},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 173, col: 15, offset: 5371},
											val:        "not",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 173, col: 21, offset: 5377},
											name: "ws",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 173, col: 27, offset: 5383},
							label: "val",
							expr: &choiceExpr{
								pos: position{line: 173, col: 32, offset: 5388},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 173, col: 32, offset: 5388},
										name: "InfixExpr",
									},
									&ruleRefExpr{
										pos:  position{line: 173, col: 44, offset: 5400},
										name: "PrefixExpr",
									},
									&ruleRefExpr{
										pos:  position{line: 173, col: 57, offset: 5413},
										name: "Term",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "InfixExpr",
			pos:  position{line: 181, col: 1, offset: 5555},
			expr: &actionExpr{
				pos: position{line: 181, col: 14, offset: 5568},
				run: (*parser).callonInfixExpr1,
				expr: &seqExpr{
					pos: position{line: 181, col: 14, offset: 5568},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 181, col: 14, offset: 5568},
							label: "left",
							expr: &ruleRefExpr{
								pos:  position{line: 181, col: 19, offset: 5573},
								name: "Term",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 181, col: 24, offset: 5578},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 181, col: 26, offset: 5580},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 181, col: 29, offset: 5583},
								name: "InfixOp",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 181, col: 37, offset: 5591},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 181, col: 39, offset: 5593},
							label: "right",
							expr: &ruleRefExpr{
								pos:  position{line: 181, col: 45, offset: 5599},
								name: "Term",
							},
						},
					},
				},
			},
		},
		{
			name: "InfixOp",
			pos:  position{line: 185, col: 1, offset: 5674},
			expr: &actionExpr{
				pos: position{line: 185, col: 12, offset: 5685},
				run: (*parser).callonInfixOp1,
				expr: &labeledExpr{
					pos:   position{line: 185, col: 12, offset: 5685},
					label: "val",
					expr: &choiceExpr{
						pos: position{line: 185, col: 17, offset: 5690},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 185, col: 17, offset: 5690},
								val:        "=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 185, col: 23, offset: 5696},
								val:        "!=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 185, col: 30, offset: 5703},
								val:        "<=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 185, col: 37, offset: 5710},
								val:        ">=",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 185, col: 44, offset: 5717},
								val:        "<",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 185, col: 50, offset: 5723},
								val:        ">",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "PrefixExpr",
			pos:  position{line: 197, col: 1, offset: 5967},
			expr: &choiceExpr{
				pos: position{line: 197, col: 15, offset: 5981},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 197, col: 15, offset: 5981},
						name: "SetEmpty",
					},
					&ruleRefExpr{
						pos:  position{line: 197, col: 26, offset: 5992},
						name: "Builtin",
					},
				},
			},
		},
		{
			name: "Builtin",
			pos:  position{line: 199, col: 1, offset: 6001},
			expr: &actionExpr{
				pos: position{line: 199, col: 12, offset: 6012},
				run: (*parser).callonBuiltin1,
				expr: &seqExpr{
					pos: position{line: 199, col: 12, offset: 6012},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 199, col: 12, offset: 6012},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 199, col: 15, offset: 6015},
								name: "Var",
							},
						},
						&litMatcher{
							pos:        position{line: 199, col: 19, offset: 6019},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 199, col: 23, offset: 6023},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 199, col: 25, offset: 6025},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 199, col: 30, offset: 6030},
								expr: &ruleRefExpr{
									pos:  position{line: 199, col: 30, offset: 6030},
									name: "Term",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 199, col: 36, offset: 6036},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 199, col: 41, offset: 6041},
								expr: &seqExpr{
									pos: position{line: 199, col: 43, offset: 6043},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 199, col: 43, offset: 6043},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 199, col: 45, offset: 6045},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 199, col: 49, offset: 6049},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 199, col: 51, offset: 6051},
											name: "Term",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 199, col: 59, offset: 6059},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 199, col: 62, offset: 6062},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Term",
			pos:  position{line: 215, col: 1, offset: 6464},
			expr: &actionExpr{
				pos: position{line: 215, col: 9, offset: 6472},
				run: (*parser).callonTerm1,
				expr: &labeledExpr{
					pos:   position{line: 215, col: 9, offset: 6472},
					label: "val",
					expr: &choiceExpr{
						pos: position{line: 215, col: 15, offset: 6478},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 215, col: 15, offset: 6478},
								name: "Comprehension",
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 31, offset: 6494},
								name: "Composite",
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 43, offset: 6506},
								name: "Scalar",
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 52, offset: 6515},
								name: "Ref",
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 58, offset: 6521},
								name: "Var",
							},
						},
					},
				},
			},
		},
		{
			name: "Comprehension",
			pos:  position{line: 219, col: 1, offset: 6552},
			expr: &ruleRefExpr{
				pos:  position{line: 219, col: 18, offset: 6569},
				name: "ArrayComprehension",
			},
		},
		{
			name: "ArrayComprehension",
			pos:  position{line: 221, col: 1, offset: 6589},
			expr: &actionExpr{
				pos: position{line: 221, col: 23, offset: 6611},
				run: (*parser).callonArrayComprehension1,
				expr: &seqExpr{
					pos: position{line: 221, col: 23, offset: 6611},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 221, col: 23, offset: 6611},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 27, offset: 6615},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 221, col: 29, offset: 6617},
							label: "term",
							expr: &ruleRefExpr{
								pos:  position{line: 221, col: 34, offset: 6622},
								name: "Term",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 39, offset: 6627},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 221, col: 41, offset: 6629},
							val:        "|",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 45, offset: 6633},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 221, col: 47, offset: 6635},
							label: "body",
							expr: &ruleRefExpr{
								pos:  position{line: 221, col: 52, offset: 6640},
								name: "Body",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 57, offset: 6645},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 221, col: 59, offset: 6647},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Composite",
			pos:  position{line: 227, col: 1, offset: 6772},
			expr: &choiceExpr{
				pos: position{line: 227, col: 14, offset: 6785},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 227, col: 14, offset: 6785},
						name: "Object",
					},
					&ruleRefExpr{
						pos:  position{line: 227, col: 23, offset: 6794},
						name: "Array",
					},
					&ruleRefExpr{
						pos:  position{line: 227, col: 31, offset: 6802},
						name: "Set",
					},
				},
			},
		},
		{
			name: "Scalar",
			pos:  position{line: 229, col: 1, offset: 6807},
			expr: &choiceExpr{
				pos: position{line: 229, col: 11, offset: 6817},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 229, col: 11, offset: 6817},
						name: "Number",
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 20, offset: 6826},
						name: "String",
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 29, offset: 6835},
						name: "Bool",
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 36, offset: 6842},
						name: "Null",
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 231, col: 1, offset: 6848},
			expr: &choiceExpr{
				pos: position{line: 231, col: 8, offset: 6855},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 231, col: 8, offset: 6855},
						name: "Scalar",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 17, offset: 6864},
						name: "Ref",
					},
					&ruleRefExpr{
						pos:  position{line: 231, col: 23, offset: 6870},
						name: "Var",
					},
				},
			},
		},
		{
			name: "Object",
			pos:  position{line: 233, col: 1, offset: 6875},
			expr: &actionExpr{
				pos: position{line: 233, col: 11, offset: 6885},
				run: (*parser).callonObject1,
				expr: &seqExpr{
					pos: position{line: 233, col: 11, offset: 6885},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 233, col: 11, offset: 6885},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 233, col: 15, offset: 6889},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 233, col: 17, offset: 6891},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 233, col: 22, offset: 6896},
								expr: &seqExpr{
									pos: position{line: 233, col: 23, offset: 6897},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 233, col: 23, offset: 6897},
											name: "Key",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 27, offset: 6901},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 233, col: 29, offset: 6903},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 33, offset: 6907},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 35, offset: 6909},
											name: "Term",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 233, col: 42, offset: 6916},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 233, col: 47, offset: 6921},
								expr: &seqExpr{
									pos: position{line: 233, col: 49, offset: 6923},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 233, col: 49, offset: 6923},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 233, col: 51, offset: 6925},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 55, offset: 6929},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 57, offset: 6931},
											name: "Key",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 61, offset: 6935},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 233, col: 63, offset: 6937},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 67, offset: 6941},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 233, col: 69, offset: 6943},
											name: "Term",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 233, col: 77, offset: 6951},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 233, col: 79, offset: 6953},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Array",
			pos:  position{line: 257, col: 1, offset: 7732},
			expr: &actionExpr{
				pos: position{line: 257, col: 10, offset: 7741},
				run: (*parser).callonArray1,
				expr: &seqExpr{
					pos: position{line: 257, col: 10, offset: 7741},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 257, col: 10, offset: 7741},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 257, col: 14, offset: 7745},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 257, col: 17, offset: 7748},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 257, col: 22, offset: 7753},
								expr: &ruleRefExpr{
									pos:  position{line: 257, col: 22, offset: 7753},
									name: "Term",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 257, col: 28, offset: 7759},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 257, col: 33, offset: 7764},
								expr: &seqExpr{
									pos: position{line: 257, col: 34, offset: 7765},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 257, col: 34, offset: 7765},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 257, col: 36, offset: 7767},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 257, col: 40, offset: 7771},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 257, col: 42, offset: 7773},
											name: "Term",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 257, col: 49, offset: 7780},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 257, col: 51, offset: 7782},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Set",
			pos:  position{line: 281, col: 1, offset: 8355},
			expr: &choiceExpr{
				pos: position{line: 281, col: 8, offset: 8362},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 281, col: 8, offset: 8362},
						name: "SetEmpty",
					},
					&ruleRefExpr{
						pos:  position{line: 281, col: 19, offset: 8373},
						name: "SetNonEmpty",
					},
				},
			},
		},
		{
			name: "SetEmpty",
			pos:  position{line: 283, col: 1, offset: 8386},
			expr: &actionExpr{
				pos: position{line: 283, col: 13, offset: 8398},
				run: (*parser).callonSetEmpty1,
				expr: &seqExpr{
					pos: position{line: 283, col: 13, offset: 8398},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 283, col: 13, offset: 8398},
							val:        "set(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 283, col: 20, offset: 8405},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 283, col: 22, offset: 8407},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SetNonEmpty",
			pos:  position{line: 289, col: 1, offset: 8495},
			expr: &actionExpr{
				pos: position{line: 289, col: 16, offset: 8510},
				run: (*parser).callonSetNonEmpty1,
				expr: &seqExpr{
					pos: position{line: 289, col: 16, offset: 8510},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 289, col: 16, offset: 8510},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 289, col: 20, offset: 8514},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 289, col: 22, offset: 8516},
							label: "head",
							expr: &ruleRefExpr{
								pos:  position{line: 289, col: 27, offset: 8521},
								name: "Term",
							},
						},
						&labeledExpr{
							pos:   position{line: 289, col: 32, offset: 8526},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 289, col: 37, offset: 8531},
								expr: &seqExpr{
									pos: position{line: 289, col: 38, offset: 8532},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 289, col: 38, offset: 8532},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 289, col: 40, offset: 8534},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 289, col: 44, offset: 8538},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 289, col: 46, offset: 8540},
											name: "Term",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 289, col: 53, offset: 8547},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 289, col: 55, offset: 8549},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Ref",
			pos:  position{line: 306, col: 1, offset: 8954},
			expr: &actionExpr{
				pos: position{line: 306, col: 8, offset: 8961},
				run: (*parser).callonRef1,
				expr: &seqExpr{
					pos: position{line: 306, col: 8, offset: 8961},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 306, col: 8, offset: 8961},
							label: "head",
							expr: &ruleRefExpr{
								pos:  position{line: 306, col: 13, offset: 8966},
								name: "Var",
							},
						},
						&labeledExpr{
							pos:   position{line: 306, col: 17, offset: 8970},
							label: "tail",
							expr: &oneOrMoreExpr{
								pos: position{line: 306, col: 22, offset: 8975},
								expr: &choiceExpr{
									pos: position{line: 306, col: 24, offset: 8977},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 306, col: 24, offset: 8977},
											name: "RefDot",
										},
										&ruleRefExpr{
											pos:  position{line: 306, col: 33, offset: 8986},
											name: "RefBracket",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "RefDot",
			pos:  position{line: 319, col: 1, offset: 9225},
			expr: &actionExpr{
				pos: position{line: 319, col: 11, offset: 9235},
				run: (*parser).callonRefDot1,
				expr: &seqExpr{
					pos: position{line: 319, col: 11, offset: 9235},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 319, col: 11, offset: 9235},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 319, col: 15, offset: 9239},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 319, col: 19, offset: 9243},
								name: "Var",
							},
						},
					},
				},
			},
		},
		{
			name: "RefBracket",
			pos:  position{line: 326, col: 1, offset: 9462},
			expr: &actionExpr{
				pos: position{line: 326, col: 15, offset: 9476},
				run: (*parser).callonRefBracket1,
				expr: &seqExpr{
					pos: position{line: 326, col: 15, offset: 9476},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 326, col: 15, offset: 9476},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 326, col: 19, offset: 9480},
							label: "val",
							expr: &choiceExpr{
								pos: position{line: 326, col: 24, offset: 9485},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 326, col: 24, offset: 9485},
										name: "Ref",
									},
									&ruleRefExpr{
										pos:  position{line: 326, col: 30, offset: 9491},
										name: "Scalar",
									},
									&ruleRefExpr{
										pos:  position{line: 326, col: 39, offset: 9500},
										name: "Var",
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 326, col: 44, offset: 9505},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Var",
			pos:  position{line: 330, col: 1, offset: 9534},
			expr: &actionExpr{
				pos: position{line: 330, col: 8, offset: 9541},
				run: (*parser).callonVar1,
				expr: &labeledExpr{
					pos:   position{line: 330, col: 8, offset: 9541},
					label: "val",
					expr: &ruleRefExpr{
						pos:  position{line: 330, col: 12, offset: 9545},
						name: "VarChecked",
					},
				},
			},
		},
		{
			name: "VarChecked",
			pos:  position{line: 335, col: 1, offset: 9667},
			expr: &seqExpr{
				pos: position{line: 335, col: 15, offset: 9681},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 335, col: 15, offset: 9681},
						label: "val",
						expr: &ruleRefExpr{
							pos:  position{line: 335, col: 19, offset: 9685},
							name: "VarUnchecked",
						},
					},
					&notCodeExpr{
						pos: position{line: 335, col: 32, offset: 9698},
						run: (*parser).callonVarChecked4,
					},
				},
			},
		},
		{
			name: "VarUnchecked",
			pos:  position{line: 339, col: 1, offset: 9763},
			expr: &actionExpr{
				pos: position{line: 339, col: 17, offset: 9779},
				run: (*parser).callonVarUnchecked1,
				expr: &seqExpr{
					pos: position{line: 339, col: 17, offset: 9779},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 339, col: 17, offset: 9779},
							name: "AsciiLetter",
						},
						&zeroOrMoreExpr{
							pos: position{line: 339, col: 29, offset: 9791},
							expr: &choiceExpr{
								pos: position{line: 339, col: 30, offset: 9792},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 339, col: 30, offset: 9792},
										name: "AsciiLetter",
									},
									&ruleRefExpr{
										pos:  position{line: 339, col: 44, offset: 9806},
										name: "DecimalDigit",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 346, col: 1, offset: 9949},
			expr: &actionExpr{
				pos: position{line: 346, col: 11, offset: 9959},
				run: (*parser).callonNumber1,
				expr: &seqExpr{
					pos: position{line: 346, col: 11, offset: 9959},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 346, col: 11, offset: 9959},
							expr: &litMatcher{
								pos:        position{line: 346, col: 11, offset: 9959},
								val:        "-",
								ignoreCase: false,
							},
						},
						&choiceExpr{
							pos: position{line: 346, col: 18, offset: 9966},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 346, col: 18, offset: 9966},
									name: "Float",
								},
								&ruleRefExpr{
									pos:  position{line: 346, col: 26, offset: 9974},
									name: "Integer",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Float",
			pos:  position{line: 359, col: 1, offset: 10365},
			expr: &choiceExpr{
				pos: position{line: 359, col: 10, offset: 10374},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 359, col: 10, offset: 10374},
						name: "ExponentFloat",
					},
					&ruleRefExpr{
						pos:  position{line: 359, col: 26, offset: 10390},
						name: "PointFloat",
					},
				},
			},
		},
		{
			name: "ExponentFloat",
			pos:  position{line: 361, col: 1, offset: 10402},
			expr: &seqExpr{
				pos: position{line: 361, col: 18, offset: 10419},
				exprs: []interface{}{
					&choiceExpr{
						pos: position{line: 361, col: 20, offset: 10421},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 361, col: 20, offset: 10421},
								name: "PointFloat",
							},
							&ruleRefExpr{
								pos:  position{line: 361, col: 33, offset: 10434},
								name: "Integer",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 361, col: 43, offset: 10444},
						name: "Exponent",
					},
				},
			},
		},
		{
			name: "PointFloat",
			pos:  position{line: 363, col: 1, offset: 10454},
			expr: &seqExpr{
				pos: position{line: 363, col: 15, offset: 10468},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 363, col: 15, offset: 10468},
						expr: &ruleRefExpr{
							pos:  position{line: 363, col: 15, offset: 10468},
							name: "Integer",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 363, col: 24, offset: 10477},
						name: "Fraction",
					},
				},
			},
		},
		{
			name: "Fraction",
			pos:  position{line: 365, col: 1, offset: 10487},
			expr: &seqExpr{
				pos: position{line: 365, col: 13, offset: 10499},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 365, col: 13, offset: 10499},
						val:        ".",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 365, col: 17, offset: 10503},
						expr: &ruleRefExpr{
							pos:  position{line: 365, col: 17, offset: 10503},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "Exponent",
			pos:  position{line: 367, col: 1, offset: 10518},
			expr: &seqExpr{
				pos: position{line: 367, col: 13, offset: 10530},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 367, col: 13, offset: 10530},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 367, col: 18, offset: 10535},
						expr: &charClassMatcher{
							pos:        position{line: 367, col: 18, offset: 10535},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 367, col: 24, offset: 10541},
						expr: &ruleRefExpr{
							pos:  position{line: 367, col: 24, offset: 10541},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 369, col: 1, offset: 10556},
			expr: &choiceExpr{
				pos: position{line: 369, col: 12, offset: 10567},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 369, col: 12, offset: 10567},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 369, col: 20, offset: 10575},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 369, col: 20, offset: 10575},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 369, col: 40, offset: 10595},
								expr: &ruleRefExpr{
									pos:  position{line: 369, col: 40, offset: 10595},
									name: "DecimalDigit",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 371, col: 1, offset: 10612},
			expr: &actionExpr{
				pos: position{line: 371, col: 11, offset: 10622},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 371, col: 11, offset: 10622},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 371, col: 11, offset: 10622},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 371, col: 15, offset: 10626},
							expr: &choiceExpr{
								pos: position{line: 371, col: 17, offset: 10628},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 371, col: 17, offset: 10628},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 371, col: 17, offset: 10628},
												expr: &ruleRefExpr{
													pos:  position{line: 371, col: 18, offset: 10629},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 371, col: 30, offset: 10641,
											},
										},
									},
									&seqExpr{
										pos: position{line: 371, col: 34, offset: 10645},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 371, col: 34, offset: 10645},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 371, col: 39, offset: 10650},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 371, col: 57, offset: 10668},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Bool",
			pos:  position{line: 380, col: 1, offset: 10926},
			expr: &choiceExpr{
				pos: position{line: 380, col: 9, offset: 10934},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 380, col: 9, offset: 10934},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 380, col: 9, offset: 10934},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 384, col: 5, offset: 11034},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 384, col: 5, offset: 11034},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 390, col: 1, offset: 11135},
			expr: &actionExpr{
				pos: position{line: 390, col: 9, offset: 11143},
				run: (*parser).callonNull1,
				expr: &litMatcher{
					pos:        position{line: 390, col: 9, offset: 11143},
					val:        "null",
					ignoreCase: false,
				},
			},
		},
		{
			name: "AsciiLetter",
			pos:  position{line: 396, col: 1, offset: 11238},
			expr: &charClassMatcher{
				pos:        position{line: 396, col: 16, offset: 11253},
				val:        "[A-Za-z_]",
				chars:      []rune{'_'},
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 398, col: 1, offset: 11264},
			expr: &charClassMatcher{
				pos:        position{line: 398, col: 16, offset: 11279},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 400, col: 1, offset: 11295},
			expr: &choiceExpr{
				pos: position{line: 400, col: 19, offset: 11313},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 400, col: 19, offset: 11313},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 400, col: 38, offset: 11332},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 402, col: 1, offset: 11347},
			expr: &charClassMatcher{
				pos:        position{line: 402, col: 21, offset: 11367},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 404, col: 1, offset: 11380},
			expr: &seqExpr{
				pos: position{line: 404, col: 18, offset: 11397},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 404, col: 18, offset: 11397},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 404, col: 22, offset: 11401},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 404, col: 31, offset: 11410},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 404, col: 40, offset: 11419},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 404, col: 49, offset: 11428},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 406, col: 1, offset: 11438},
			expr: &charClassMatcher{
				pos:        position{line: 406, col: 17, offset: 11454},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 408, col: 1, offset: 11461},
			expr: &charClassMatcher{
				pos:        position{line: 408, col: 24, offset: 11484},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 410, col: 1, offset: 11491},
			expr: &charClassMatcher{
				pos:        position{line: 410, col: 13, offset: 11503},
				val:        "[0-9a-f]",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name:        "ws",
			displayName: "\"whitespace\"",
			pos:         position{line: 412, col: 1, offset: 11513},
			expr: &oneOrMoreExpr{
				pos: position{line: 412, col: 20, offset: 11532},
				expr: &charClassMatcher{
					pos:        position{line: 412, col: 20, offset: 11532},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 414, col: 1, offset: 11544},
			expr: &zeroOrMoreExpr{
				pos: position{line: 414, col: 19, offset: 11562},
				expr: &choiceExpr{
					pos: position{line: 414, col: 21, offset: 11564},
					alternatives: []interface{}{
						&charClassMatcher{
							pos:        position{line: 414, col: 21, offset: 11564},
							val:        "[ \\t\\r\\n]",
							chars:      []rune{' ', '\t', '\r', '\n'},
							ignoreCase: false,
							inverted:   false,
						},
						&ruleRefExpr{
							pos:  position{line: 414, col: 33, offset: 11576},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 416, col: 1, offset: 11588},
			expr: &actionExpr{
				pos: position{line: 416, col: 12, offset: 11599},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 416, col: 12, offset: 11599},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 416, col: 12, offset: 11599},
							expr: &charClassMatcher{
								pos:        position{line: 416, col: 12, offset: 11599},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 416, col: 19, offset: 11606},
							val:        "#",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 416, col: 23, offset: 11610},
							label: "text",
							expr: &zeroOrMoreExpr{
								pos: position{line: 416, col: 28, offset: 11615},
								expr: &charClassMatcher{
									pos:        position{line: 416, col: 28, offset: 11615},
									val:        "[^\\r\\n]",
									chars:      []rune{'\r', '\n'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 422, col: 1, offset: 11750},
			expr: &notExpr{
				pos: position{line: 422, col: 8, offset: 11757},
				expr: &anyMatcher{
					line: 422, col: 9, offset: 11758,
				},
			},
		},
	},
}

func (c *current) onProgram1(vals interface{}) (interface{}, error) {
	var buf []interface{}

	if vals == nil {
		return buf, nil
	}

	ifaceSlice := vals.([]interface{})
	head := ifaceSlice[0]
	buf = append(buf, head)
	for _, tail := range ifaceSlice[1].([]interface{}) {
		stmt := tail.([]interface{})[1]
		buf = append(buf, stmt)
	}

	return buf, nil
}

func (p *parser) callonProgram1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onProgram1(stack["vals"])
}

func (c *current) onStmt1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonStmt1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStmt1(stack["val"])
}

func (c *current) onParseError1() (interface{}, error) {
	panic(fmt.Sprintf("no match found, unexpected '%s'", c.text))
}

func (p *parser) callonParseError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParseError1()
}

func (c *current) onPackage1(val interface{}) (interface{}, error) {
	// All packages are implicitly declared under the default root document.
	path := RefTerm(DefaultRootDocument)
	switch v := val.(*Term).Value.(type) {
	case Ref:
		// Convert head of package Ref to String because it will be prefixed
		// with the root document variable.
		head := v[0]
		head = StringTerm(string(head.Value.(Var)))
		head.Location = v[0].Location
		tail := v[1:]
		if !tail.IsGround() {
			return nil, fmt.Errorf("package name cannot contain variables: %v", v)
		}

		// We do not allow non-string values in package names.
		// Because documents are typically represented as JSON, non-string keys are
		// not allowed for now.
		// TODO(tsandall): consider special syntax for namespacing under arrays.
		for _, p := range tail {
			_, ok := p.Value.(String)
			if !ok {
				return nil, fmt.Errorf("package name cannot contain non-string values: %v", v)
			}
		}
		path.Value = append(path.Value.(Ref), head)
		path.Value = append(path.Value.(Ref), tail...)
	case Var:
		s := StringTerm(string(v))
		s.Location = val.(*Term).Location
		path.Value = append(path.Value.(Ref), s)
	}
	pkg := &Package{Location: currentLocation(c), Path: path.Value.(Ref)}
	return pkg, nil
}

func (p *parser) callonPackage1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPackage1(stack["val"])
}

func (c *current) onImport1(path, alias interface{}) (interface{}, error) {
	imp := &Import{}
	imp.Location = currentLocation(c)
	imp.Path = path.(*Term)
	if err := IsValidImportPath(imp.Path.Value); err != nil {
		return nil, err
	}
	if alias == nil {
		return imp, nil
	}
	aliasSlice := alias.([]interface{})
	// Import definition above describes the "alias" slice. We only care about the "Var" element.
	imp.Alias = aliasSlice[3].(*Term).Value.(Var)
	return imp, nil
}

func (p *parser) callonImport1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport1(stack["path"], stack["alias"])
}

func (c *current) onRule1(name, key, value, body interface{}) (interface{}, error) {

	rule := &Rule{}
	rule.Location = currentLocation(c)
	rule.Name = name.(*Term).Value.(Var)

	if key != nil {
		keySlice := key.([]interface{})
		// Rule definition above describes the "key" slice. We care about the "Term" element.
		rule.Key = keySlice[3].(*Term)

		var closure interface{}
		WalkClosures(rule.Key, func(x interface{}) bool {
			closure = x
			return true
		})

		if closure != nil {
			return nil, fmt.Errorf("head cannot contain closures (%v appears in key)", closure)
		}
	}

	if value != nil {
		valueSlice := value.([]interface{})
		// Rule definition above describes the "value" slice. We care about the "Term" element.
		rule.Value = valueSlice[len(valueSlice)-1].(*Term)

		var closure interface{}
		WalkClosures(rule.Value, func(x interface{}) bool {
			closure = x
			return true
		})

		if closure != nil {
			return nil, fmt.Errorf("head cannot contain closures (%v appears in value)", closure)
		}
	}

	if key == nil && value == nil {
		rule.Value = BooleanTerm(true)
	}

	if key != nil && value != nil {
		switch rule.Key.Value.(type) {
		case Var, String, Ref: // nop
		default:
			return nil, fmt.Errorf("head of object rule must have string, var, or ref key (%s is not allowed)", rule.Key)
		}
	}

	// Rule definition above describes the "body" slice. We only care about the "Body" element.
	rule.Body = body.([]interface{})[3].(Body)

	return rule, nil
}

func (p *parser) callonRule1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRule1(stack["name"], stack["key"], stack["value"], stack["body"])
}

func (c *current) onBody1(head, tail interface{}) (interface{}, error) {
	var buf Body
	buf = append(buf, head.(*Expr))
	for _, s := range tail.([]interface{}) {
		expr := s.([]interface{})[3].(*Expr)
		buf = append(buf, expr)
	}
	return buf, nil
}

func (p *parser) callonBody1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBody1(stack["head"], stack["tail"])
}

func (c *current) onExpr1(neg, val interface{}) (interface{}, error) {
	expr := &Expr{}
	expr.Location = currentLocation(c)
	expr.Negated = neg != nil
	expr.Terms = val
	return expr, nil
}

func (p *parser) callonExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpr1(stack["neg"], stack["val"])
}

func (c *current) onInfixExpr1(left, op, right interface{}) (interface{}, error) {
	return []*Term{op.(*Term), left.(*Term), right.(*Term)}, nil
}

func (p *parser) callonInfixExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInfixExpr1(stack["left"], stack["op"], stack["right"])
}

func (c *current) onInfixOp1(val interface{}) (interface{}, error) {
	op := string(c.text)
	for _, b := range Builtins {
		if string(b.Infix) == op {
			op = string(b.Name)
		}
	}
	operator := VarTerm(op)
	operator.Location = currentLocation(c)
	return operator, nil
}

func (p *parser) callonInfixOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInfixOp1(stack["val"])
}

func (c *current) onBuiltin1(op, head, tail interface{}) (interface{}, error) {
	buf := []*Term{op.(*Term)}
	if head == nil {
		return buf, nil
	}
	buf = append(buf, head.(*Term))

	// PrefixExpr above describes the "tail" structure. We only care about the "Term" elements.
	tailSlice := tail.([]interface{})
	for _, v := range tailSlice {
		s := v.([]interface{})
		buf = append(buf, s[len(s)-1].(*Term))
	}
	return buf, nil
}

func (p *parser) callonBuiltin1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBuiltin1(stack["op"], stack["head"], stack["tail"])
}

func (c *current) onTerm1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonTerm1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTerm1(stack["val"])
}

func (c *current) onArrayComprehension1(term, body interface{}) (interface{}, error) {
	ac := ArrayComprehensionTerm(term.(*Term), body.(Body))
	ac.Location = currentLocation(c)
	return ac, nil
}

func (p *parser) callonArrayComprehension1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArrayComprehension1(stack["term"], stack["body"])
}

func (c *current) onObject1(head, tail interface{}) (interface{}, error) {
	obj := ObjectTerm()
	obj.Location = currentLocation(c)

	// Empty object.
	if head == nil {
		return obj, nil
	}

	// Object definition above describes the "head" structure. We only care about the "Key" and "Term" elements.
	headSlice := head.([]interface{})
	obj.Value = append(obj.Value.(Object), Item(headSlice[0].(*Term), headSlice[len(headSlice)-1].(*Term)))

	// Non-empty object, remaining key/value pairs.
	tailSlice := tail.([]interface{})
	for _, v := range tailSlice {
		s := v.([]interface{})
		// Object definition above describes the "tail" structure. We only care about the "Key" and "Term" elements.
		obj.Value = append(obj.Value.(Object), Item(s[3].(*Term), s[len(s)-1].(*Term)))
	}

	return obj, nil
}

func (p *parser) callonObject1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObject1(stack["head"], stack["tail"])
}

func (c *current) onArray1(head, tail interface{}) (interface{}, error) {

	arr := ArrayTerm()
	arr.Location = currentLocation(c)

	// Empty array.
	if head == nil {
		return arr, nil
	}

	// Non-empty array, first element.
	arr.Value = append(arr.Value.(Array), head.(*Term))

	// Non-empty array, remaining elements.
	tailSlice := tail.([]interface{})
	for _, v := range tailSlice {
		s := v.([]interface{})
		// Array definition above describes the "tail" structure. We only care about the "Term" elements.
		arr.Value = append(arr.Value.(Array), s[len(s)-1].(*Term))
	}

	return arr, nil
}

func (p *parser) callonArray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArray1(stack["head"], stack["tail"])
}

func (c *current) onSetEmpty1() (interface{}, error) {
	set := SetTerm()
	set.Location = currentLocation(c)
	return set, nil
}

func (p *parser) callonSetEmpty1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetEmpty1()
}

func (c *current) onSetNonEmpty1(head, tail interface{}) (interface{}, error) {
	set := SetTerm()
	set.Location = currentLocation(c)

	val := set.Value.(*Set)
	val.Add(head.(*Term))

	tailSlice := tail.([]interface{})
	for _, v := range tailSlice {
		s := v.([]interface{})
		// SetNonEmpty definition above describes the "tail" structure. We only care about the "Term" elements.
		val.Add(s[len(s)-1].(*Term))
	}

	return set, nil
}

func (p *parser) callonSetNonEmpty1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetNonEmpty1(stack["head"], stack["tail"])
}

func (c *current) onRef1(head, tail interface{}) (interface{}, error) {

	ref := RefTerm(head.(*Term))
	ref.Location = currentLocation(c)

	tailSlice := tail.([]interface{})
	for _, v := range tailSlice {
		ref.Value = append(ref.Value.(Ref), v.(*Term))
	}

	return ref, nil
}

func (p *parser) callonRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRef1(stack["head"], stack["tail"])
}

func (c *current) onRefDot1(val interface{}) (interface{}, error) {
	// Convert the Var into a string because 'foo.bar.baz' is equivalent to 'foo["bar"]["baz"]'.
	str := StringTerm(string(val.(*Term).Value.(Var)))
	str.Location = currentLocation(c)
	return str, nil
}

func (p *parser) callonRefDot1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRefDot1(stack["val"])
}

func (c *current) onRefBracket1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonRefBracket1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRefBracket1(stack["val"])
}

func (c *current) onVar1(val interface{}) (interface{}, error) {
	return val.([]interface{})[0], nil
}

func (p *parser) callonVar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVar1(stack["val"])
}

func (c *current) onVarChecked4(val interface{}) (bool, error) {
	return IsKeyword(string(val.(*Term).Value.(Var))), nil
}

func (p *parser) callonVarChecked4() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVarChecked4(stack["val"])
}

func (c *current) onVarUnchecked1() (interface{}, error) {
	str := string(c.text)
	variable := VarTerm(str)
	variable.Location = currentLocation(c)
	return variable, nil
}

func (p *parser) callonVarUnchecked1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVarUnchecked1()
}

func (c *current) onNumber1() (interface{}, error) {
	f, ok := new(big.Float).SetString(string(c.text))
	if !ok {
		// This indicates the grammar is out-of-sync with what the string
		// representation of floating point numbers. This should not be
		// possible.
		panic("illegal value")
	}
	num := NumberTerm(json.Number(f.String()))
	num.Location = currentLocation(c)
	return num, nil
}

func (p *parser) callonNumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumber1()
}

func (c *current) onString1() (interface{}, error) {
	// TODO : the forward slash (solidus) is not a valid escape in Go, it will
	// fail if there's one in the string
	v, err := strconv.Unquote(string(c.text))
	str := StringTerm(v)
	str.Location = currentLocation(c)
	return str, err
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onBool2() (interface{}, error) {
	bol := BooleanTerm(true)
	bol.Location = currentLocation(c)
	return bol, nil
}

func (p *parser) callonBool2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool2()
}

func (c *current) onBool4() (interface{}, error) {
	bol := BooleanTerm(false)
	bol.Location = currentLocation(c)
	return bol, nil
}

func (p *parser) callonBool4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool4()
}

func (c *current) onNull1() (interface{}, error) {
	null := NullTerm()
	null.Location = currentLocation(c)
	return null, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
}

func (c *current) onComment1(text interface{}) (interface{}, error) {
	comment := NewComment(ifaceSliceToByteSlice(text))
	comment.Location = currentLocation(c)
	return comment, nil
}

func (p *parser) callonComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onComment1(stack["text"])
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}