package gcd

// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm

// function extended_gcd(a, b)
//     s := 0;    old_s := 1
//     t := 1;    old_t := 0
//     r := b;    old_r := a
//
//     while r ≠ 0
//         quotient := old_r div r
//         (old_r, r) := (r, old_r - quotient * r)
//         (old_s, s) := (s, old_s - quotient * s)
//         (old_t, t) := (t, old_t - quotient * t)
//
//     output "Bézout coefficients:", (old_s, old_t)
//     output "greatest common divisor:", old_r
//     output "quotients by the gcd:", (t, s)
//

// X는 Extended Euclidean algorithm을 계산합니다. 다음 불변식을 만족합니다.
//   a*s + b*t = gcd(a, b) = r
func X(a, b int) (s int, t int, r int) { return xgcd(a, b) }

// E는 Euclidean algorithm을 계산합니다. 다음 불변식을 만족합니다.
//   gcd(a, b) = r
func E(a, b int) (r int) { return gcd(a, b) }

// IsCoprime은 두 수가 서로소인지 검사합니다.
func IsCoprime(a, b int) bool { return isCoprime(a, b) }

// MultiplicativeInverse은 두 수가 서로소일 때 a의 역원을 계산합니다.
func MultiplicativeInverse(a, b int) int { return multiplicativeInverse(a, b) }

// a*s + b*t = gcd(a, b) = r
func xgcd(a, b int) (s int, t int, r int) {
	s0, s1 := 1, 0
	t0, t1 := 0, 1
	r0, r1 := a, b
	for r1 != 0 {
		quotient := r0 / r1
		r0, r1 = r1, r0-quotient*r1
		s0, s1 = s1, s0-quotient*s1
		t0, t1 = t1, t0-quotient*t1
	}
	return s0, t0, r0
}

func gcd(a, b int) int {
	r0, r1 := a, b
	for r1 != 0 {
		quotient := r0 / r1
		r0, r1 = r1, r0-quotient*r1
	}
	return r0
}

func isCoprime(a, b int) bool {
	return gcd(a, b) == 1
}

// a*s = 1 (mod b)
func multiplicativeInverse(a, b int) (s int) {
	s, _, _ = xgcd(a, b)
	return s
}
