def load(filename: str) -> list[tuple[str]]:
    with open(filename) as f:
        return list(map(str.split, f))

def main():
    codegen = load('codegen.txt')
    generics = load('generics.txt')
    print(f"{'Test name timed in ns/op':<50}Codegen   \tGenerics  \tRegression [%]")
    print("===========================================================================================")
    for c, g in zip(codegen, generics):
        if c == g or c[0] == g[0] == 'go':
            continue
        tc, tg = float(c[2]), float(g[2])
        print(f"{c[0]:<50}{tc:<10}\t{tg:<10}\t{100*tg/tc-100:.2f}")
        #print(c[0], tc, tg, round(tg/tc-1, 2), sep='\t')

if __name__ == '__main__':
    main()
