const MAXBITS: usize = 4;

#[derive(Debug)]
struct Node {
    value: usize,
    left: Option<Box<Node>>,
    right: Option<Box<Node>>,
}

impl Node {
    fn new(value: usize) -> Node {
        Node {
            value,
            left: None,
            right: None,
        }
    }

    fn insert(&mut self, value: usize) {
        if bits(value, 0, 1) == 1 {
            if let Some(ref mut right) = self.right {
                right.insert(value >> 1);
            } else {
                self.right = Some(Box::new(Node::new(1)));
            }
            return;
        }
        if let Some(ref mut left) = self.left {
            left.insert(value >> 1);
        } else {
            self.left = Some(Box::new(Node::new(0)));
        }
    }

    fn search(&self, value: usize) -> bool {
        if bits(value, 0, 1) == 1 {
            return match self.right {
                Some(ref right) => right.search(value >> 1),
                None => value == 1,
            };
        }
        match self.left {
            Some(ref left) => left.search(value >> 1),
            None => value == 0,
        }
    }
}

fn bits(n: usize, shift: usize, size: usize) -> usize {
    (n >> shift) & !(!0 << size)
}

struct Tree {
    root: Option<Box<Node>>,
}

impl Tree {
    fn new() -> Tree {
        Tree { root: None }
    }

    fn insert(&mut self, value: usize) {
        match self.root {
            Some(ref mut root) => root.insert(value),
            None => self.root = Some(Box::new(Node::new(value))),
        }
    }

    fn search(&self, value: usize) -> bool {
        match self.root {
            Some(ref root) => root.search(value),
            None => false,
        }
    }
}

fn main() {
    let mut tree = Tree::new();
    for i in 0..8 {
        println!("insert {}", i);
        tree.insert(i);
    }
    println!("{:#?}", tree.root);
    for i in 0..10 {
        if tree.search(i) {
            println!("found {}", i);
        } else {
            println!("not found {}", i);
        }
    }
}
