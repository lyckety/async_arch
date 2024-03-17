package consumer

import "testing"

func Test_parseTaskDescription(t *testing.T) {
	type args struct {
		desc string
	}
	tests := []struct {
		name       string
		args       args
		wantJiraID string
		wantTitle  string
	}{
		{
			name: "1",
			args: args{
				desc: "[UBERPOPUG-1] My Task 1",
			},
			wantJiraID: "UBERPOPUG-1",
			wantTitle:  "My Task 1",
		},
		{
			name: "2",
			args: args{
				desc: "My Task 2",
			},
			wantJiraID: "",
			wantTitle:  "My Task 2",
		},
		{
			name: "3",
			args: args{
				desc: "[]My Task 3",
			},
			wantJiraID: "",
			wantTitle:  "My Task 3",
		},
		{
			name: "4",
			args: args{
				desc: "[]My Task 4",
			},
			wantJiraID: "",
			wantTitle:  "My Task 4",
		},
		{
			name: "5",
			args: args{
				desc: "",
			},
			wantJiraID: "",
			wantTitle:  "",
		},
		{
			name: "6",
			args: args{
				desc: "[UBERPOPUG-6]",
			},
			wantJiraID: "UBERPOPUG-6",
			wantTitle:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJiraID, gotTitle := parseTaskDescription(tt.args.desc)
			if gotJiraID != tt.wantJiraID {
				t.Errorf("parseTaskDescription() gotJiraID = %v, wantJiraID %v", gotJiraID, tt.wantJiraID)
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("parseTaskDescription() gotTitle = %v, wantTitle %v", gotTitle, tt.wantTitle)
			}
		})
	}
}
