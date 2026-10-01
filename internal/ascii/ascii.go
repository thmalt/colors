package ascii

func EqualFold(s, t string) bool {
	for i, n := 0, min(len(s), len(t)); i < n; i++ {
		sr := s[i]
		tr := t[i]

		if tr == sr {
			continue
		}

		if tr < sr {
			tr, sr = sr, tr
		}

		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			continue
		}

		return false
	}

	return len(s) == len(t)
}

func ToLowerBytes(b []byte) {
	for i := range b {
		c := b[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
			b[i] = c
		}
	}
}
