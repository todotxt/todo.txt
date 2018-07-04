package main

import (
	"fmt"

	parsec "github.com/prataprc/goparsec"
)

func main() {
	ast := parsec.NewAST("TODO", 1000)

	tag := ast.OrdChoice("TAG", nil,
		parsec.Token(`[+][[:graph:]]+`, "PLUSTAG"),
		parsec.Token(`[@][[:graph:]]+`, "ATTAG"),
		parsec.Token(`[#][[:graph:]]+`, "HASHTAG"),
	)

	kvPair := ast.And("KVPAIR", nil,
		// we use this long regexp instead of :graph: to omit ':' There's probably a regexpy
		// way of doing this better
		parsec.Token(`[^#@+:][A-Za-z0-9!"#$%&'()*+,\-./;<=>?@[\\\]^_{|}~]*`, "KEY"),
		parsec.TokenExact(":", "COLON"),
		parsec.TokenExact(`[A-Za-z0-9!"#$%&'()*+,\-./;<=>?@[\\\]^_{|}~]+`, "VALUE"),
	)

	word := ast.OrdChoice("WORD", nil,
		parsec.Token(`[^@#+][[:graph:]]*`, "WORD"),
		parsec.Token(`[@#+]$`, "WORD"),
	)

	token := ast.OrdChoice("TOKEN", nil, kvPair, word, tag)

	createdDate := parsec.Token(`[0-9]{4}-[0-9]{2}-[0-9]{2}`, "CREATIONDATE")
	completeDate := parsec.Token(`[0-9]{4}-[0-9]{2}-[0-9]{2}`, "COMPLETIONDATE")

	priority := parsec.Token(`\([A-Z]\)[[:space:]]+`, "PRIORITY")

	completeMark := parsec.Token(`x[[:space:]]+`, "COMPLETED")

	TODO := ast.And("TODO", nil,
		ast.OrdChoice("PREAMBLE", nil,
			ast.And("PREAMBLE", nil,
				completeMark,
				ast.Maybe("PRIORITY", nil, priority),
				ast.Maybe("COMPLETEDAT", nil, completeDate),
				ast.Maybe("CREATEDAT", nil, createdDate),
			),
			ast.And("PREAMBLE", nil,
				ast.Maybe("PRIORITY", nil, priority),
				ast.Maybe("CREATEDAT", nil, createdDate),
			),
		),
		ast.Many("WORDS", nil, token),
	)

	examples := []string{
		"x (A) 2018-07-30 2018-07-31 Some todo item with +ProjectTag @atTag #hashtag and key:value pair",
		"(A) Thank Mom for the meatballs @phone",
		"(B) Schedule Goodwill pickup +GarageSale @phone",
		"Post signs around the neighborhood +GarageSale",
		"@GroceryStore Eskimo pies",
		"(A) Thank Mom for the meatballs @phone",
		"(B) Schedule Goodwill pickup +GarageSale @phone",
		"(B) Schedule Goodwill pickup +GarageSale @phone",
		"Post signs around the neighborhood +GarageSale",
		"Really gotta call Mom (A) @phone @someday",
		"(b) Get back to the boss",
		"(B)->Submit TPS report",
		"2011-03-02 Document +TodoTxt task format", // This is not a completion date because it's not complete
		"(A) 2011-03-02 Call Mom",
		"(A) Call Mom 2011-03-02",
		"(A) Call Mom +Family +PeaceLoveAndHappiness @iphone @phone",
		"Email SoAndSo at soandso@example.com",
		"Learn how to add 2+2",
		"x 2011-03-03 Call Mom", // this is a completion date
		"xylophone lesson",
		"X 2012-01-01 Make resolutions",
		"(A) x Find ticket prices",
		"x 2011-03-02 2011-03-01 Review Tim's pull request +TodoTxtTouch @github",
		"Some example with key:value and due:2010-01-02",
	}

	for i, todo := range examples {
		println()
		fmt.Printf("[%v]: %s\n", i, todo)
		scanner := parsec.NewScanner([]byte(todo))
		ast.Parsewith(TODO, scanner)
		ast.Prettyprint()
		ast.Reset()
	}
}
