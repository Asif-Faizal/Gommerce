-- Migration: Create users table
-- Description: Initial migration to create the users table with necessary fields
-- Author: System
-- Date: 2024-03-19

-- Create users table with all required fields
CREATE TABLE IF NOT EXISTS users (
    -- Primary key field
    id INT UNSIGNED AUTO_INCREMENT,
    
    -- User's personal information
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    
    -- Timestamps for record keeping
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Primary key constraint
    PRIMARY KEY (id),
    -- Ensure email uniqueness
    UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;