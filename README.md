algorithm-class
===============

이 저장소는 인천대학교 컴퓨터공학과 2019년 2학기 채진석 교수의 알고리즘 수업
시간에 구현한 여러 알고리즘의 소스 코드를 포함하고 있습니다. 또한 정렬
알고리즘과 트리 알고리즘에 대하여 성늘 비교를 위한 성능 그래프를 그릴 수
있습니다. 아래 코드는 실제로 제출한 코드입니다. 해당 수업에서 A+을 받았습니다.  

모든 소스 코드는 Go로 작성되었습니다. 소스 코드를 개선시킬 방법이 있다면 이슈를
통해 또는 메일로 알려주시면 감사하겠습니다. 계정 아이디에 @gmail.com을 붙이면
계정 주소입니다.

최소한 아래 알고리즘을 수업에서 공부하게 됩니다. 이 저장소는 아래 알고리즘을
구현하고 있습니다.

- Sorting
  - Bubble Sort
  - Cocktail Shaker Sort
  - Cracking Quick Sort Algorithm
  - Exchange Sort
  - Heap Sort
  - Merge Sort
  - Natural Merge Sort
  - Natural Merge Sort (Heap)
  - Quick Sort
    - Insertion Quick Sort
    - Median Quick Sort
    - Median Insertion Quick Sort
  - Selection Sort
  - Shell Sort
  - Tournament Sort
- Binary Search
- Eratosthenes Sieve
- String Searching
  - BoyerMoore (BM) String Search (Bad char skip applied only)
  - KMP String Search
  - Naive String Search
  - Rabin-Karp String Search
  - Simple Pattern Matching
- Tree
  - AVL Tree
  - Binary Tree
  - Digital Search Tree
  - Huffman Encoding
  - Patricia Tree
  - Radix Tree
  - Red-Black Tree
- Geometry
  - Closest Two Points
  - Graham Scan
  - Package Wrapping
  - Point-In-Polygon Detection
  - Segment Crossing Detection

트리 알고리즘은 삽입 연산만 구현되어 있습니다.  격자 교점 찾기 등 일부
알고리즘에 대해서 배웠으나 구현하지 않은 알고리즘이 일부 있습니다.

# 성능 측정
정렬 알고리즘과 트리 알고리즘에 대해서 성능을 측정할 수 있는 방법을 제공합니다.
`sort`와 `search`를 단순히 빌드하여 알고리즘 간 성능을 측정할 수 있습니다. `go
build` 명령를 통해 간단하게 빌드할 수 있습니다.

아래 그래프는 토너먼트 정렬에서 비교 횟수에 대한 그래프입니다.

![tournament-sort - fuzz input](https://user-images.githubusercontent.com/56159921/87287499-243ae700-c535-11ea-87b3-ec66e0b716de.jpeg)

아래 그래프는 역순으로 정렬된 입력에 대하여 AVL 트리와 적흑트리의 노드 접근
횟수에 대한 그래프입니다.

![AVL, RB, - reversed-input](https://user-images.githubusercontent.com/56159921/87287494-2309ba00-c535-11ea-9ec4-075a6fbe3c86.jpeg)

아래 그래프는 역순으로 정렬된 입력에 대하여 63 비트를 사용하는 디지털 탐색
트리와 패트리샤 트리 그리고 이진 트리의 노드 접근 횟수에 대한 그래프입니다.

![DST, PT, BT, - reversed input](https://user-images.githubusercontent.com/56159921/87287498-243ae700-c535-11ea-8f84-a94367aba2dd.jpeg)


공개된 소스 코드는 학습을 목적으로 자유롭게 사용하실 수 있습니다. 소스 코드를
사용한다는 것은 소스 코드 사용으로 인한 모든 책임을 감수함을 의미합니다.

소스 코드를 사용한다면 반드시 출처를 남겨주시기 바랍니다. 출처를 표시할 때
저장소 주소인 https://github.com/gopherclass/algorithm-class 를 반드시 포함해야
합니다.

