This fixture starts a single workflow that adds a lot of event histories `(LengthOfHistory)` and then fails `(WillFailOrNot)`.

Our UI should handle this properly by showing that this workflow has failed, despite the failure event being on the next page.

Used in:
- https://github.com/temporalio/web/issues/300
