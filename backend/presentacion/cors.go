package presentacion

import "net/http"

// conCORS aplica la política de origen; la autenticación se configura aparte.
func conCORS(siguiente http.Handler, origen string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		if r.Header.Get("Origin") == origen {
			w.Header().Set("Access-Control-Allow-Origin", origen)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			if r.Header.Get("Origin") != origen {
				responder(w, http.StatusForbidden, map[string]string{"error": "Origen no permitido."})
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Header.Get("Origin") != "" && r.Header.Get("Origin") != origen {
			responder(w, http.StatusForbidden, map[string]string{"error": "Origen no permitido."})
			return
		}
		siguiente.ServeHTTP(w, r)
	})
}
