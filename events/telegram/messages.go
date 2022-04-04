package telegram

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
)

const msgHelp = `Commands: 
/meal - return a random meal
/add - add new meal to the collection, see /add_sample.
	Format:

	/add

	photo: <photo url*>
	===
	name: <meal name*>
	===
	instructions: <instructions paragraph>
	===
	description: <description paragraph>
	
	- Replace '<>' with content.
	- Delimeter === has to be provided after each section except of the last one.
	- * means the value is required.
	- The order of the sections can be changed.
`

const msgHello = `Hi there! 👾
Ramadan has come and you don't know what to cook? 
Just ask me!

Send /meal command and I will find a meal for you.

For more information and other commands send /help`

const msgAddSample = `/add

photo: https://loopbarbados.com/sites/default/files/styles/blog_image_style/public/blogimages/Bajan-Backed-Chicken-Recipe.jpg?itok=bt3ZknC8
===
name: BAKED BAJAN CHICKEN
===
instructions:
1. Preheat oven to 450 degrees.
2. Smear the Bajan seasoning inside and under the skin of the chicken.
3. etc...
===
description: Delicious baked chicked with an ancient history cooked in Eastern Countries.`
