# ecommerce

This is a sample demo that depicts a set of ecommerce applications built with a microservices paradigm.

## Usecases

We'll target the following usecases for an MVP version of the demo. More usecases can be added later, but we'll focus on these as a start.

In general, we will support 4 basic services:
* Catalog: Services to view all of the items available in the store catalog.
* Users: Services for a user to be able to login, and access their profiles
* Cart: Services for a user to add specific items from the Catalog to a cart.
* Checkout: Services that allows a user to checkout with 

### Must Haves:
#### User
* A user should be able to login with a username and password.
* The act of logging in should return a token that can be used later on to validate the user session.
* A user should be able to retrieve his profile.
* A user consists of a firstname, lastname, address, id, email and password.

#### Catalog
* Catalog will represent items available in the ecommerce site.
* Each item will have a name, price, SKU (unique ID).

#### Cart
* Catalog items can be added to a cart for a logged in user
* Cart will display all currently added items for a specific user.

#### Checkout
* Checkout will take the items in a cart, and create an order for the user.

### Nice to Haves:
#### User:
* A user can create an account.

#### Catalog:
* Ability to add a list and discount prices

#### Cart
* A user can save their carts

#### Checkout
* tbd

### API's

#### User
* POST /users/ (Login, returns a session token)
* GET /users/<userID> (Authenticated, returns user information)

#### Catalog
* GET /catalog (returns all catalog info)
* GET /catalog/<itemID> (returns a specific item in catalog based on ID)

#### Cart
* GET /cart/<cartID> (Authenticated, returns a )
* POST /cart (Authenticated, creates or updates a cart for a user)

#### Checkout
* POST /checkout?cart=<cartID> (Authenticated, creates an order for the cart ID passed)