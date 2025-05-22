package helper

import "time"

// GetLocation retorna a localização padrão do sistema (America/Bahia)
func GetLocation() *time.Location {
	loc, err := time.LoadLocation("America/Bahia")
	if err != nil {
		// Fallback para UTC se não conseguir carregar a timezone
		return time.UTC
	}
	return loc
}

// TimeNow retorna o tempo atual na timezone America/Bahia
func TimeNow() time.Time {
	return time.Now().In(GetLocation())
}
