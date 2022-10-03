**ADR**

1. Use gorm (ORM for golang)
  
To interact with database application use GORM.

GORM allow to use soft delete functionality. But it's not necessary in this solution.
So I decide to not using it. This make solution simpler.

**API**

*Book API*

*Author API*
 
- Add author
- Delete author
  - soft delete
- Edit author
- List authors
  - in default should return existing authors
  - could return soft deleted authors
- Get author
  - in default should return existing author or null
  - could return soft delete author

TODO:
- create controller for author API
- add necessary endpoint to controller
- implement every endpoint
- unit tests
- integration tests

