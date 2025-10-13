# SOLUTION FOR MORE IDOMATIC GO

I have structured the steps I took into separate commits, to tackle each instructions of the [ASSIGNMENT.md](ASSIGNMENT.md), starting with the idiomatic golang file structure - with a Domain friendly approach - and ending with last reverse engineering of the functions in `app/api/response.go` by following thier unit test.

## Refactoring the Project into a more idiomatic Go structure

The first issue (as the Assignmet points out) is that the `models.ProductsRepository` is heavily coupled with the `catalog.CatalogHandler`, which is renders the code less maitainable and less scalable. 

To solve this, I used an approach that focuses on the Domain of the _Product_ , thus creating a packgae and a folder named `product` , and to avoid complicating the project's structure into a Hexagonal with multiple nested layers (inra, app and domain), I've opten to encapsulate all realtes components: _Repository, Handler, Use Case ..._ under the same package, while focusing on breaking the existing coupling of _Repository_ and _Http Handler_ by introducing a Domain _Service_ that will implement the use cases and intermediates both infra components (_Respository_ and _Handler_)

And in turn, moving all of these folders and code into the usual Go's `internal` folder.

## Adding the Project Model

At first glance I was going to add the _Category_ entity alongside the _Product_ one, however, since one of the assignments indicated a lifecycle of categories indepenent from the _Product_ one, with it's own endpoint to handle listing and creation, I decided to add new domain folder/package called `category`.

When it come to the migration sql scripts, I add new ones instead of modifying exiting one to reflect the changes made to the schema and new link between `product` and `category`.

## Support offset pagination and filtering

Mostly refactoring the existing components:
- the `handler`'s GET call to parse the new prams: _offset, limit, category_ and _priceLessThan_ 
- the _Serice_ has the Category `repository` now to convert the category _code_ into category _ID_ (Assuming that the front/client of this API will use category codes instead of free string text of the category), also the _Service_ validates the parameters received from the `handler` (Via a new DTO `FindCriteria`) 
- the Product `repository` has a more dynamic query based on the if the filter params (_categoryId_ or _priceLessThn_) are present, with a side query to get the total Products based on the filters.
- the Category `repository` implements the `GetCategoryByCode` so the _Serivce_ might use it to convert the Code to ID (as previously mentioned). Also in preparation for the next task of the assignent (to handle the `/category` endpoint)

## Adding Categories endpoint

For listing all, I've implemented the `GET /categories`, making it simple without pagination, filtering nor total. As for the creation I followed the standard of using `POST /category` with a No Content HTTP status 204 as a response.
The implementation follows the previous structure and component pattern made in the `product` domain: _handler, service_ and _repository_

Also added unit tests for both endpoints/handlers

