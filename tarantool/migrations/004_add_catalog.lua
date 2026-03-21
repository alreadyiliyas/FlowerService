return {
    id = "004_add_catalog",
    up = function ()
        local c = box.schema.space.create("categories", { if_not_exists = true })
        c:format({
            {name="id",   type="unsigned"},
            {name="name", type="string"},
            {name="slug", type="string"},
            {name="description", type="string", is_nullable=true},
            {name="image_url",   type="string", is_nullable=true},
            {name="created_at",  type="unsigned"},
            {name="updated_at",  type="unsigned", is_nullable=true},
        })
        log.info("--------------------------> create categories <-------------------------->")
        box.schema.sequence.create("categories_seq", { if_not_exists = true } )
        c:create_index("primary", { parts={{field="id", type="unsigned"}}, sequence="categories_seq", if_not_exists=true })
        c:create_index("name", { parts={{field="name", type="string"}}, unique=true, if_not_exists=true })
        c:create_index("slug", { parts={{field="slug", type="string"}}, unique=true, if_not_exists=true })

        log.info("--------------------------> create categories_seq <-------------------------->")
        local p = box.schema.space.create("products", { if_not_exists = true })
        p:format({
            {name="id",             type="unsigned"},
            {name="name",           type="string"},
            {name="slug",           type="string"},
            {name="description",    type="string", is_nullable=true},
            {name="category_id",    type="unsigned"},
            {name="seller_id",      type="unsigned"},
            {name="is_available",   type="boolean"},
            {name="currency",       type="string"},
            {name="main_image_url", type="string"},
            {name="images",         type="array"},
            {name="sizes",          type="array"},
            {name="price_per_stem", type="unsigned"},
            {name="min_stems",      type="unsigned"},
            {name="max_stems",      type="unsigned"},
            {name="composition",    type="array"},
            {name="discount",       type="map", is_nullable=true},
            {name="version",        type="unsigned"},
            {name="created_at",     type="unsigned"},
            {name="updated_at",     type="unsigned", is_nullable=true},
            {name="deleted_at",     type="unsigned", is_nullable=true},
        })
        log.info("--------------------------> create products <-------------------------->")
        
        box.schema.sequence.create("products_seq", { if_not_exists = true } )
        p:create_index("primary", { parts={{field="id", type="unsigned"}}, sequence="products_seq", if_not_exists=true })
        p:create_index("slug", { parts={{field="slug", type="string"}}, unique=true, if_not_exists=true })
        p:create_index("category_id", { parts={{field="category_id", type="unsigned"}}, unique=false, if_not_exists=true })
        p:create_index("seller_id", { parts={{field="seller_id", type="unsigned"}}, unique=false, if_not_exists=true })
        p:create_index("is_available", { parts={{field="is_available", type="boolean"}}, unique=false, if_not_exists=true })
        p:create_index("created_at", { parts={{field="created_at", type="unsigned"}}, unique=false, if_not_exists=true })
        log.info("--------------------------> create products_seq <-------------------------->")
    end
}
