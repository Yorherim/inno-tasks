package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type args struct {
	mtx [][]int
	ua  []int
}
type test struct {
	name string
	args args
	want int
}

func TestEvalSequence(t *testing.T) {

	mtx1 := [][]int{
		{0, 2, 3, 0, 0},
		{2, 0, 0, 1, 1},
		{3, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0},
	}

	tests := []test{
		{
			name: "mtx 5 verticals 100%",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 1, 0, 2},
			},
			want: 100,
		},
		{
			name: "mtx 5 verticals 0%",
			args: args{
				mtx: mtx1,
				ua:  []int{},
			},
			want: 0,
		},
		{
			name: "mtx 5 verticals 50%",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 1, 0},
			},
			want: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvalSequence(tt.args.mtx, tt.args.ua)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalcUserGrade(t *testing.T) {
	matrix := [][]int{
		{0, 2, 3, 0, 0},
		{2, 0, 0, 1, 1},
		{3, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0},
	}

	tests := []test{
		{
			name: "6",
			args: args{
				mtx: matrix,
				ua:  []int{4, 1, 0, 2},
			},
			want: 6,
		},
		{
			name: "0",
			args: args{
				mtx: matrix,
				ua:  []int{},
			},
			want: 0,
		},
		{
			name: "3",
			args: args{
				mtx: matrix,
				ua:  []int{4, 1, 0},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcUserGrade(tt.args.mtx, tt.args.ua)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidation(t *testing.T) {
	mtx1 := [][]int{
		{0, 2, 3, 0, 0},
		{2, 0, 0, 1, 1},
		{3, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0},
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "error nil",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 1, 0, 2},
			},
			want: nil,
		},
		{
			name: "empty matrix",
			args: args{
				mtx: [][]int{},
				ua:  []int{4, 1, 0, 2},
			},
			want: errors.New(ErrorEmptyMatrix),
		},
		{
			name: "empty user answers",
			args: args{
				mtx: mtx1,
				ua:  []int{},
			},
			want: errors.New(ErrorUserAnswer),
		},
		{
			name: "matrix not square",
			args: args{
				mtx: [][]int{
					{0, 2, 3, 0, 0},
					{2, 0, 0, 1, 1},
					{3, 0, 0, 0, 0},
					{0, 1, 0, 0, 0},
					{0, 1, 0, 0},
				},
				ua: []int{4, 1, 0, 2},
			},
			want: errors.New(ErrorMatrixNotSquare),
		},
		{
			name: "graph loop",
			args: args{
				mtx: [][]int{
					{0, 2, 3, 0, 0},
					{2, 0, 0, 1, 1},
					{3, 0, 1, 0, 0},
					{0, 1, 0, 0, 0},
					{0, 1, 0, 0, 0},
				},
				ua: []int{4, 1, 0, 2},
			},
			want: errors.New(ErrorGraphLoop),
		},
		{
			name: "user answers range",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 1, 0, 2, 1, 1, 1, 1, 1, 1, 1},
			},
			want: errors.New(ErrorUserAnswersRange),
		},
		{
			name: "user answers incorrect",
			args: args{
				mtx: mtx1,
				ua:  []int{10, 1, 0, 2},
			},
			want: errors.New(ErrorUserAnswersIncorrect),
		},
		{
			name: "user answer not unique",
			args: args{
				mtx: mtx1,
				ua:  []int{4, 4, 0, 2},
			},
			want: errors.New(ErrorUserAnswerNotUnique),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation(tt.args.mtx, tt.args.ua)
			if tt.want == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.Error())
			}
		})
	}
}
