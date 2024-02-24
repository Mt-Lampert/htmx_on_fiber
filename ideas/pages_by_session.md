
# Pages by session

## Ingredients

- altogether 22 existing contacts

## Description

The idea is a paging feature:

1. The landing page shows a total of 10 contacts
0. Below is a button `More`
0. If the button is clicked, a total of 15 contacts is shown and the focus will
   be laid to the bottom of the table.
0. If the button is clicked again, a total of 20 contacts is shown and the
   focus will be laid to the bottom of the table again.
0. If all contacts are being shown on the page, the `More` button is being replaced
   by a `Less` button.
0. 





### Thoughts about the implementation.

1. The SQL for getting all contacts must be updated by a `LIMIT ?` clause
0. The current `LIMIT` numer must be saved. For the moment a global variable in 
   the `main` package will do the job (“Simplest solution wins!”)
   - since there is only one user, advanced techniques (see below) won't be
      necessary at the moment.
   - integrated solution: we initialize the set number to 2 and add `LIMIT ?`
     to the SQL. 
        - then we pass `sets * numOfContacts` to the database function at startup.
        - when the user clicks `More`, we raise `sets` by 1 and repeat the query.
        - in the handler we check how many contacts there are in total in the
          database. if `sets * numOfContacts` is greater or equal to the total,
          we switch the button 



### YAGNI Ideas

YAGNI stands for “You ain’t gonna need it!”

- [Fiber sessions](https://docs.gofiber.io/api/middleware/session): Only for
  more than one user.

