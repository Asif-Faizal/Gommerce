-- Migration: Drop users table
-- Description: Reverts the creation of the users table
-- Author: System
-- Date: 2024-03-19

-- Drop the users table if it exists
DROP TABLE IF EXISTS users;