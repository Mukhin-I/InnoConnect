import CategoryButton from './CategoryButton';
import './CategorySelector.css';

const CategorySelector = ({
    categories = [],
    selectedCategory,
    onSelectCategory,
    className = '',
    ...props
}) => {
    return (
        <div className={`category-selector ${className}`} {...props}>
            {categories.map((category) => (
                <CategoryButton
                    key={category.id}
                    icon={selectedCategory === category.id ? category.activeIcon : category.inactiveIcon}
                    label={category.label}
                    isActive={selectedCategory === category.id}
                    onClick={() => onSelectCategory(category.id)}
                />
            ))}
        </div>
    );
};

export default CategorySelector;