# Risks and Technical Debt

## Risk: Data Inconsistency in Borrowing

**Problem**  
Book borrowing operations update the database in multiple steps without transactions. Under concurrent usage, this can cause:

- 📉 **Inventory mismatches**: Books show as available when they're borrowed  
- 👥 **User assignment errors**: Loans disappear or get assigned to wrong users  
- 🔢 **Count inaccuracies**: Total borrowed books don't match user records  

**Required Fix**  
Add database transaction support to make all updates atomic.