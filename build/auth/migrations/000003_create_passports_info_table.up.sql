CREATE TABLE IF NOT EXISTS passports_info(
    id SERIAL UNIQUE,
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    passp_series VARCHAR(4),
    Passp_number VARCHAR(6),
	passp_date_of_issue DATE,    
	passp_department_code VARCHAR(6),
	passp_issue_by VARCHAR(255),
	passp_name VARCHAR(63) NOT NULL,  
	passp_lastName VARCHAR(63) NOT NULL,   
	passp_patronymic VARCHAR(63),
	passp_sex VARCHAR(15) NOT NULL,    
	passp_date_of_birth DATE, 
	passp_place_of_birth VARCHAR(127),
	passp_registration VARCHAR(127)
);