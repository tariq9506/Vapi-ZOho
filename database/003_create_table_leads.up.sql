 CREATE TABLE leads (
    id VARCHAR(50) NOT NULL UNIQUE, 
    full_name VARCHAR(255),         
    lead_status VARCHAR(50),        
    phone VARCHAR(15),              
    created_at TIMESTAMP,           
    district VARCHAR(255),         
    gender VARCHAR(10),             
    query_param VARCHAR(100),      
    session_type VARCHAR(50),
    call_ended_reasone VARCHAR(150),
    dialing_code VARCHAR(3),
    PRIMARY KEY (id),              
);