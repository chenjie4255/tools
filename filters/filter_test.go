package filters

import "testing"

func Test_filterString(t *testing.T) {
	type args struct {
		sources []string
		filter  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "正常取值(成功)",
			args: args{[]string{"cn", "us"}, "us"},
			want: true,
		},
		{
			name: "正常取值(失败)",
			args: args{[]string{"cn", "us"}, "hk"},
			want: false,
		},
		{
			name: "空值取值",
			args: args{[]string{}, "us"},
			want: true,
		},
		{
			name: "过滤器空值取值",
			args: args{[]string{"jp"}, ""},
			want: true,
		},
		{
			name: "反向取值(成功)",
			args: args{[]string{"!jp"}, "cn"},
			want: true,
		},
		{
			name: "反向取值(失败)",
			args: args{[]string{"!jp"}, "JP"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterString(tt.args.sources, tt.args.filter); got != tt.want {
				t.Errorf("filterString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterRanges(t *testing.T) {
	type args struct {
		ranges []string
		val    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "+范围取值(成功)",
			args: args{[]string{"1.2.0+"}, "1.3.0"},
			want: true,
		},
		{
			name: "+范围取值(成功2)",
			args: args{[]string{"1.2.0+"}, "1.2.0"},
			want: true,
		},
		{
			name: "+范围取值(失败)",
			args: args{[]string{"1.2.0+"}, "1.1.0"},
			want: false,
		},
		{
			name: "-范围取值(成功)",
			args: args{[]string{"1.2.0-"}, "1.1.0"},
			want: true,
		},
		{
			name: "-范围取值(失败)",
			args: args{[]string{"1.2.0-"}, "1.2.1"},
			want: false,
		},
		{
			name: "区间范围取值(成功)",
			args: args{[]string{"1.2.0-1.3.0"}, "1.2.1"},
			want: true,
		},
		{
			name: "区间范围取值(失败)",
			args: args{[]string{"1.2.0-1.3.0"}, "1.1.1"},
			want: false,
		},
		{
			name: "直接匹配(成功)",
			args: args{[]string{"1.2.0"}, "1.2.0"},
			want: true,
		},
		{
			name: "直接匹配(失败)",
			args: args{[]string{"1.2.0.1"}, "1.2.0"},
			want: false,
		},
		{
			name: "线上修复",
			args: args{[]string{"3.0.0-"}, "3.0.1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckVersionInRanges(tt.args.ranges, tt.args.val); got != tt.want {
				t.Errorf("filterRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
