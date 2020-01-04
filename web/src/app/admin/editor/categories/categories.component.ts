import { Component, OnInit } from '@angular/core';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

@Component({
  selector: 'app-post-categories-editor',
  templateUrl: './categories.component.html',
  styleUrls: ['./categories.component.scss'],
})
export class PostCategoriesEditorComponent extends AbstractPostEditorComponent implements OnInit {

  /**
   * List of categories to-be rendered as sidebar menu-item(s)
   */
  categories: Array<Category>;

  ngOnInit(): void {
    this.apollo.query({
      query: gql`
        {
          categories { name slug }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ categories: Array<Category> }>): Array<Category> => result.data.categories),
    ).subscribe((categories: Array<Category>): void => {
      this.categories = categories;
    })
  }

}
