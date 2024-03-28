CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT
        NOT NULL
        UNIQUE
        CONSTRAINT name_length CHECK (char_length(username) <= 255),
    password_hash TEXT
        NOT NULL
        CONSTRAINT password_hash_length CHECK (char_length(password_hash) <= 511),
    create_time TIMESTAMP
        NOT NULL,
    image_path TEXT DEFAULT ('default.jpg')
        NOT NULL
        CONSTRAINT image_path_length CHECK (char_length(image_path) <= 255)
    
);
CREATE TABLE IF NOT EXISTS market (
    id UUID PRIMARY KEY,
    title TEXT
        NOT NULL
        CONSTRAINT title_length CHECK (char_length(title) <= 255),
    description TEXT
        NOT NULL
        CONSTRAINT description_length CHECK (char_length(description) <= 1000),
    create_time TIMESTAMP
        NOT NULL,
    image_path TEXT DEFAULT ('default.jpg')
        NOT NULL
        CONSTRAINT image_path_length CHECK (char_length(image_path) <= 255),
    price NUMERIC 
        NOT NULL
        CONSTRAINT price_constrain CHECK (price>0),
    owner UUID REFERENCES users (id)
    
);