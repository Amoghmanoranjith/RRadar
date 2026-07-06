## Entities:
**Post**
   1. title
   2. link
   3. content
   4. author 
   5. published time
   6. subreddit

**ClassifiedPost**
1. Title
2. Link
3. Content
4. Author
5. Published time
6. subreddit
7. Interesting
8. Confidence
9. Reason

## DB Entity:
**Post**:
1. subreddit
2. author
3. published


## Flow:
1. Scraper scrapes subreddits and returns a list of []byte
2. An adapter converts this to a list of Posts
3. Filter the posts using the db
4. Update the db with the latest post for respective subreddit
5. Give the filtered list of posts to a classifier
6. Classifier goes through each post and creates respective ClassifiedPost
7. The list of classifiedPosts is given to a notifier
8. Notifier notifies mediums based on interesting nature