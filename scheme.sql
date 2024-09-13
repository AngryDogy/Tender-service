CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS employee (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        username VARCHAR(50) UNIQUE NOT NULL,
                                        first_name VARCHAR(50),
                                        last_name VARCHAR(50),
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE IF NOT EXISTS organization (
                                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                            name VARCHAR(100) NOT NULL,
                                            description TEXT,
                                            type organization_type,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
                                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                                        organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
                                                        user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TYPE service_type as ENUM(
    'Construction',
    'Delivery',
    'Manufacture'
    );

CREATE TYPE status as ENUM(
    'Created',
    'Published',
    'Closed'
    'Canceled'
    );

CREATE TYPE author_type as ENUM (
    'Organization',
    'User'
);

CREATE TABLE IF NOT EXISTS tender (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                      name VARCHAR(100) NOT NULL,
                                      description TEXT,
                                      service_type service_type,
                                      status status DEFAULT 'Created',
                                      organization_id UUID REFERENCES organization(id) on delete cascade ,
                                      creator_username VARCHAR(100) REFERENCES employee(username) on delete CASCADE ,
                                      version integer DEFAULT 1,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS tender_version(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID REFERENCES tender(id),
    version integer,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type service_type
);


CREATE TABLE IF NOT EXISTS bid(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status status DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) on delete cascade,
    author_type author_type,
    author_id UUID,
    version integer DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid_version(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id),
    version integer,
    name VARCHAR(100),
    description TEXT
);

CREATE TABLE IF NOT EXISTS feedback(
    id UUID PRIMARY KEY  default uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) on delete cascade ,
    description TEXT,
    username VARCHAR(100),
    created_at TIMESTAMP default current_timestamp
);


INSERT INTO employee(id, username, first_name, last_name) VALUES('00155f0f-e5b3-47c1-8ff9-fba3566401b0', 'angrydog', 'Vlad', 'Hober');

INSERT INTO employee(id, username, first_name, last_name) VALUES('25b7f208-2a54-4313-a166-c1a8658a00cc', 'repkin', 'Jenya', 'Repkin');

INSERT INTO organization(id, name, description, type) VALUES('00155f0f-e5b3-47c1-8ff9-fba3566401b0', 'ITMO', 'the best uni', 'IE');

INSERT INTO organization(id, name, description, type) VALUES('e0e483f4-ccb7-4ea9-af47-15772368ba7b', 'petrgu', 'petrozavodsk', 'IE');

INSERT INTO organization_responsible(id, organization_id, user_id) VALUES('cf543bc3-5532-4d29-9051-168be56b26c7', '00155f0f-e5b3-47c1-8ff9-fba3566401b0', '00155f0f-e5b3-47c1-8ff9-fba3566401b0');

INSERT INTO organization_responsible(id, organization_id, user_id) VALUES ( 'd51a1595-5451-46a6-9688-37bf5d1be59d', 'e0e483f4-ccb7-4ea9-af47-15772368ba7b', '25b7f208-2a54-4313-a166-c1a8658a00cc')