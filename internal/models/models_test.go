package models

import "testing"

func TestProductInputValidate(t *testing.T) {
	tests := []struct {
		name    string
		in      ProductInput
		wantErr bool
	}{
		{"válido", ProductInput{Name: "Tornillo", SKU: "TOR-001", Quantity: 10}, false},
		{"sin nombre", ProductInput{SKU: "TOR-001"}, true},
		{"nombre solo espacios", ProductInput{Name: "   ", SKU: "TOR-001"}, true},
		{"sin sku", ProductInput{Name: "Tornillo"}, true},
		{"cantidad negativa", ProductInput{Name: "Tornillo", SKU: "TOR-001", Quantity: -1}, true},
		{"cantidad cero", ProductInput{Name: "Tornillo", SKU: "TOR-001", Quantity: 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovementInputValidate(t *testing.T) {
	tests := []struct {
		name     string
		in       MovementInput
		wantErr  bool
		wantType string
	}{
		{"entrada válida", MovementInput{ProductID: 1, Type: "IN", Quantity: 5}, false, "IN"},
		{"normaliza minúsculas", MovementInput{ProductID: 1, Type: " out ", Quantity: 5}, false, "OUT"},
		{"tipo inválido", MovementInput{ProductID: 1, Type: "MOVE", Quantity: 5}, true, ""},
		{"cantidad cero", MovementInput{ProductID: 1, Type: "IN", Quantity: 0}, true, ""},
		{"cantidad negativa", MovementInput{ProductID: 1, Type: "OUT", Quantity: -3}, true, ""},
		{"sin producto", MovementInput{Type: "IN", Quantity: 1}, true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.in.Type != tt.wantType {
				t.Fatalf("Type = %q, want %q", tt.in.Type, tt.wantType)
			}
		})
	}
}
