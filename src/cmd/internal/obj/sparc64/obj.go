// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
	"log"
)

// TODO(aram):
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {
	cursym.Text.Pc = 0
	cursym.Args = cursym.Text.To.Val.(int32)
	cursym.Locals = int32(cursym.Text.To.Offset)

	// For future use by oplook and friends.
	for p := cursym.Text; p != nil; p = p.Link {
		p.From.Class = aclass(&p.From)
		if p.From3 != nil {
			p.From3.Class = aclass(p.From3)
		}
		p.To.Class = aclass(&p.To)
	}

	// Find leaf subroutines,
	// Strip NOPs.
	var q *obj.Prog
	var q1 *obj.Prog
	for p := cursym.Text; p != nil; p = p.Link {
		switch p.As {
		case obj.ATEXT:
			p.Mark |= LEAF

		case obj.ARET:
			break

		case obj.ANOP:
			q1 = p.Link
			q.Link = q1 /* q is non-nop */
			q1.Mark |= p.Mark
			continue

		case obj.AJMP, AFBA,
			obj.ADUFFZERO,
			obj.ADUFFCOPY:
			cursym.Text.Mark &^= LEAF
			fallthrough

		case ABN, ABNE, ABE, ABG, ABLE, ABGE, ABL, ABGU, ABLEU, ABCC, ABCS, ABPOS, ABNEG, ABVC, ABVS,
			ABRZ, ABRLEZ, ABRLZ, ABRNZ, ABRGZ, ABRGEZ,
			AFBN, AFBU, AFBG, AFBUG, AFBL, AFBUL, AFBLG, AFBNE, AFBE, AFBUE, AFBGE, AFBUGE, AFBLE, AFBULE, AFBO:
			q1 = p.Pcond

			if q1 != nil {
				for q1.As == obj.ANOP {
					q1 = q1.Link
					p.Pcond = q1
				}
			}

			break
		}

		q = p
	}

	for p := cursym.Text; p != nil; p = p.Link {
		switch p.As {
		case obj.ATEXT:
			if cursym.Text.Mark&LEAF != 0 {
				cursym.Leaf = 1
			}
		}
	}
}

func relinv(a int) int {
	switch a {
	case obj.AJMP:
		return ABN
	case ABN:
		return obj.AJMP
	case ABE:
		return ABNE
	case ABNE:
		return ABE
	case ABG:
		return ABLE
	case ABLE:
		return ABG
	case ABGE:
		return ABL
	case ABL:
		return ABGE
	case ABGU:
		return ABLEU
	case ABLEU:
		return ABGU
	case ABCC:
		return ABCS
	case ABCS:
		return ABCC
	case ABPOS:
		return ABNEG
	case ABNEG:
		return ABPOS
	case ABVC:
		return ABVS
	case ABVS:
		return ABVC
	}

	log.Fatalf("unknown relation: %s", Anames[a])
	return 0
}

var unaryDst = map[int]bool{
	obj.ACALL: true,
	AWORD:     true,
	ADWORD:    true,
}

var Linksparc64 = obj.LinkArch{
	ByteOrder:  binary.BigEndian,
	Name:       "sparc64",
	Thechar:    'u',
	Preprocess: preprocess,
	Assemble:   span,
	Follow:     follow,
	UnaryDst:   unaryDst,
	Minlc:      4,
	Ptrsize:    8,
	Regsize:    8,
}
